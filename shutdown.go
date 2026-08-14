package main

import (
	"context"
	"sort"
	"sync"
	"time"

	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

// DefaultShutdownPhase is the phase assigned to components registered via
// RegisterComponent. Components sharing a phase shut down concurrently;
// phases run sequentially in ascending order. A caller that never assigns an
// explicit phase gets every component in DefaultShutdownPhase, which
// reproduces the pre-phase behavior of shutting down everything at once.
//
// Invariant: the traffic-drain phase (shutdownPhaseDrainTraffic in main.go)
// must be strictly earlier than DefaultShutdownPhase. Otherwise an unphased
// component would land in the same phase as the drain, or later ones, and
// could be torn down concurrently with in-flight requests instead of after
// they drain - the exact hazard phased shutdown exists to prevent.
const DefaultShutdownPhase = 0

// ShutdownManager manages graceful shutdown for all components
type ShutdownManager struct {
	ctx          context.Context
	cancel       context.CancelFunc
	components   []ShutdownComponent
	wg           sync.WaitGroup
	shutdownOnce sync.Once
	mu           sync.Mutex
}

// ShutdownComponent represents a component that needs graceful shutdown.
// Phase controls ordering relative to other components: lower phases shut
// down first, and shutdown only moves to the next phase once every
// component in the current phase has returned or the overall shutdown
// deadline has passed.
type ShutdownComponent struct {
	Shutdown func(context.Context) error
	Name     string
	Phase    int
}

// NewShutdownManager creates a new shutdown manager
func NewShutdownManager(ctx context.Context) *ShutdownManager {
	ctx, cancel := context.WithCancel(ctx)
	return &ShutdownManager{
		ctx:    ctx,
		cancel: cancel,
	}
}

// RegisterComponent registers a component for graceful shutdown in
// DefaultShutdownPhase. Use RegisterComponentWithPhase to sequence a
// component before or after others (for example, draining traffic before
// closing the backends that in-flight requests depend on).
func (sm *ShutdownManager) RegisterComponent(name string, shutdown func(context.Context) error) {
	sm.RegisterComponentWithPhase(name, DefaultShutdownPhase, shutdown)
}

// RegisterComponentWithPhase registers a component for graceful shutdown in
// the given phase. Shutdown runs phases sequentially in ascending order;
// components within the same phase shut down concurrently.
func (sm *ShutdownManager) RegisterComponentWithPhase(name string, phase int, shutdown func(context.Context) error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.components = append(sm.components, ShutdownComponent{
		Name:     name,
		Phase:    phase,
		Shutdown: shutdown,
	})
}

// RunGoroutine starts a goroutine that respects the shutdown context
func (sm *ShutdownManager) RunGoroutine(name string, fn func(context.Context)) {
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		cfgMutex.RLock()
		logger := cfg.Logger
		cfgMutex.RUnlock()
		if logger != nil {
			logger.Debug(&libpack_logging.LogMessage{
				Message: "Starting managed goroutine",
				Pairs:   map[string]any{"name": name},
			})
		}
		fn(sm.ctx)
		cfgMutex.RLock()
		logger = cfg.Logger
		cfgMutex.RUnlock()
		if logger != nil {
			logger.Debug(&libpack_logging.LogMessage{
				Message: "Managed goroutine finished",
				Pairs:   map[string]any{"name": name},
			})
		}
	}()
}

// Shutdown initiates graceful shutdown of all components
func (sm *ShutdownManager) Shutdown(timeout time.Duration) error {
	var err error
	sm.shutdownOnce.Do(func() {
		err = sm.doShutdown(timeout)
	})
	return err
}

// doShutdown performs the actual shutdown logic
func (sm *ShutdownManager) doShutdown(timeout time.Duration) error {
	cfgMutex.RLock()
	logger := cfg.Logger
	cfgMutex.RUnlock()
	if logger != nil {
		logger.Info(&libpack_logging.LogMessage{
			Message: "Initiating graceful shutdown",
		})
	}

	// Cancel the context to signal all goroutines to stop
	sm.cancel()

	// Create a timeout context for component shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	// Shutdown all registered components, grouped by phase. Phases run
	// sequentially in ascending order against the single shutdownCtx budget
	// above (no per-phase timeout reset), so a component in an earlier phase
	// only shuts down once the whole shutdown has enough of that shared
	// budget remaining. This budget bounds the phases loop and the
	// componentsDone wait below - it does not bound all of Shutdown(): the
	// separate goroutinesDone wait further down (for RunGoroutine-managed
	// goroutines) uses its own fresh time.After(timeout) budget, so
	// worst-case total shutdown time is up to ~2x timeout. That is
	// pre-existing, intentional behavior, unchanged here.
	sm.mu.Lock()
	components := make([]ShutdownComponent, len(sm.components))
	copy(components, sm.components)
	sm.mu.Unlock()

	phaseGroups := make(map[int][]ShutdownComponent, len(components))
	phases := make([]int, 0, len(components))
	for _, comp := range components {
		if _, seen := phaseGroups[comp.Phase]; !seen {
			phases = append(phases, comp.Phase)
		}
		phaseGroups[comp.Phase] = append(phaseGroups[comp.Phase], comp)
	}
	sort.Ints(phases)

	// allWg tracks every component across every phase, so the final wait
	// below reports success only once nothing is left running anywhere.
	var allWg sync.WaitGroup
	for _, phase := range phases {
		var phaseWg sync.WaitGroup
		for _, comp := range phaseGroups[phase] {
			phaseWg.Add(1)
			allWg.Add(1)
			go func(c ShutdownComponent) {
				defer phaseWg.Done()
				defer allWg.Done()
				cfgMutex.RLock()
				logger := cfg.Logger
				cfgMutex.RUnlock()
				if logger != nil {
					logger.Info(&libpack_logging.LogMessage{
						Message: "Shutting down component",
						Pairs:   map[string]any{"component": c.Name, "phase": c.Phase},
					})
				}
				if err := c.Shutdown(shutdownCtx); err != nil {
					cfgMutex.RLock()
					logger := cfg.Logger
					cfgMutex.RUnlock()
					if logger != nil {
						logger.Error(&libpack_logging.LogMessage{
							Message: "Error shutting down component",
							Pairs: map[string]any{
								"component": c.Name,
								"phase":     c.Phase,
								"error":     err.Error(),
							},
						})
					}
				}
			}(comp)
		}

		phaseDone := make(chan struct{})
		go func() {
			phaseWg.Wait()
			close(phaseDone)
		}()

		// Bound the wait for this phase by the overall shutdown budget, not
		// a fresh per-phase timeout. If the deadline fires while this phase
		// is still running, its stragglers keep running in the background
		// (they are never killed, only ctx-cancelled) but we move on so the
		// next phase still gets attempted instead of deadlocking here. If it
		// is specifically the drain phase that times out, later phases go
		// on to close backends while requests still draining in the
		// background depend on them - an unavoidable consequence of
		// enforcing a fixed shutdown budget.
		select {
		case <-phaseDone:
		case <-shutdownCtx.Done():
			cfgMutex.RLock()
			logger := cfg.Logger
			cfgMutex.RUnlock()
			if logger != nil {
				logger.Warning(&libpack_logging.LogMessage{
					Message: "Shutdown phase timed out, proceeding to next phase",
					Pairs:   map[string]any{"phase": phase},
				})
			}
		}
	}

	// Wait for all components, across all phases, to finish shutting down.
	componentsDone := make(chan struct{})
	go func() {
		allWg.Wait()
		close(componentsDone)
	}()

	// Wait for goroutines with timeout
	goroutinesDone := make(chan struct{})
	go func() {
		sm.wg.Wait()
		close(goroutinesDone)
	}()

	select {
	case <-componentsDone:
		cfgMutex.RLock()
		logger := cfg.Logger
		cfgMutex.RUnlock()
		if logger != nil {
			logger.Info(&libpack_logging.LogMessage{
				Message: "All components shut down successfully",
			})
		}
	case <-shutdownCtx.Done():
		cfgMutex.RLock()
		logger := cfg.Logger
		cfgMutex.RUnlock()
		if logger != nil {
			logger.Warning(&libpack_logging.LogMessage{
				Message: "Component shutdown timed out",
			})
		}
	}

	select {
	case <-goroutinesDone:
		cfgMutex.RLock()
		logger := cfg.Logger
		cfgMutex.RUnlock()
		if logger != nil {
			logger.Info(&libpack_logging.LogMessage{
				Message: "All goroutines finished",
			})
		}
	case <-time.After(timeout):
		cfgMutex.RLock()
		logger := cfg.Logger
		cfgMutex.RUnlock()
		if logger != nil {
			logger.Warning(&libpack_logging.LogMessage{
				Message: "Some goroutines didn't finish within timeout",
			})
		}
	}

	return nil
}
