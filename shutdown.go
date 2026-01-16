package main

import (
	"context"
	"sync"
	"time"

	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

// ShutdownManager manages graceful shutdown for all components
type ShutdownManager struct {
	ctx          context.Context
	cancel       context.CancelFunc
	components   []ShutdownComponent
	wg           sync.WaitGroup
	shutdownOnce sync.Once
	mu           sync.Mutex
}

// ShutdownComponent represents a component that needs graceful shutdown
type ShutdownComponent struct {
	Shutdown func(context.Context) error
	Name     string
}

// NewShutdownManager creates a new shutdown manager
func NewShutdownManager(ctx context.Context) *ShutdownManager {
	ctx, cancel := context.WithCancel(ctx)
	return &ShutdownManager{
		ctx:    ctx,
		cancel: cancel,
	}
}

// RegisterComponent registers a component for graceful shutdown
func (sm *ShutdownManager) RegisterComponent(name string, shutdown func(context.Context) error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.components = append(sm.components, ShutdownComponent{
		Name:     name,
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

	// Shutdown all registered components
	sm.mu.Lock()
	components := make([]ShutdownComponent, len(sm.components))
	copy(components, sm.components)
	sm.mu.Unlock()

	var shutdownWg sync.WaitGroup
	for _, comp := range components {
		shutdownWg.Add(1)
		go func(c ShutdownComponent) {
			defer shutdownWg.Done()
			cfgMutex.RLock()
			logger := cfg.Logger
			cfgMutex.RUnlock()
			if logger != nil {
				logger.Info(&libpack_logging.LogMessage{
					Message: "Shutting down component",
					Pairs:   map[string]any{"component": c.Name},
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
							"error":     err.Error(),
						},
					})
				}
			}
		}(comp)
	}

	// Wait for all components to shutdown
	componentsDone := make(chan struct{})
	go func() {
		shutdownWg.Wait()
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
