package main

import (
	"context"
	"testing"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
)

// Pins theorems/Theorems/EventCleaner.lean: a retention below one day would delete all event history.
func TestEnableHasuraEventCleaner_RejectsNonPositiveRetention(t *testing.T) {
	for _, days := range []int{0, -1, -30} {
		cfgMutex.Lock()
		cfg = &config{Logger: libpack_logger.New()}
		cfg.HasuraEventCleaner.Enable = true
		cfg.HasuraEventCleaner.EventMetadataDb = "postgres://fake:fake@127.0.0.1:1/fake"
		cfg.HasuraEventCleaner.ClearOlderThan = days
		cfgMutex.Unlock()

		ctx, cancel := context.WithCancel(context.Background())
		err := enableHasuraEventCleaner(ctx)
		cancel()
		assert.Error(t, err, "days=%d", days)
	}
}
