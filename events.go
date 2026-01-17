package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

const (
	initialDelay    = 60 * time.Second
	cleanupInterval = 1 * time.Hour
)

// Use parameterized queries to prevent SQL injection
// Cast $1 to interval type to allow proper parameterized interval values
var delQueries = [...]string{
	"DELETE FROM hdb_catalog.event_invocation_logs WHERE created_at < NOW() - $1::INTERVAL",
	"DELETE FROM hdb_catalog.event_log WHERE created_at < NOW() - $1::INTERVAL",
	"DELETE FROM hdb_catalog.hdb_action_log WHERE created_at < NOW() - $1::INTERVAL",
	"DELETE FROM hdb_catalog.hdb_cron_event_invocation_logs WHERE created_at < NOW() - $1::INTERVAL",
	"DELETE FROM hdb_catalog.hdb_scheduled_event_invocation_logs WHERE created_at < NOW() - $1::INTERVAL",
}

func enableHasuraEventCleaner(ctx context.Context) error {
	cfgMutex.RLock()
	if !cfg.HasuraEventCleaner.Enable {
		cfgMutex.RUnlock()
		return nil
	}

	eventMetadataDb := cfg.HasuraEventCleaner.EventMetadataDb
	if eventMetadataDb == "" {
		logger := cfg.Logger
		cfgMutex.RUnlock()

		logger.Warning(&libpack_logger.LogMessage{
			Message: "Event metadata db URL not specified, event cleaner not active",
		})
		return nil
	}

	clearOlderThan := cfg.HasuraEventCleaner.ClearOlderThan
	logger := cfg.Logger
	cfgMutex.RUnlock()

	logger.Info(&libpack_logger.LogMessage{
		Message: "Event cleaner enabled",
		Pairs:   map[string]any{"interval_in_days": clearOlderThan},
	})

	// Parse pool configuration
	poolConfig, err := pgxpool.ParseConfig(eventMetadataDb)
	if err != nil {
		return err
	}

	// Set connection pool limits
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to create connection pool",
			Pairs:   map[string]any{"error": err.Error()},
		})
		return err
	}

	go func() {
		defer pool.Close()

		// Wait for initial delay or context cancellation
		select {
		case <-ctx.Done():
			return
		case <-time.After(initialDelay):
		}

		logger.Info(&libpack_logger.LogMessage{
			Message: "Initial cleanup of old events",
		})
		cleanEvents(ctx, pool, clearOlderThan, logger)

		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Info(&libpack_logger.LogMessage{
					Message: "Stopping event cleaner",
				})
				return
			case <-ticker.C:
				logger.Info(&libpack_logger.LogMessage{
					Message: "Cleaning up old events",
				})
				cleanEvents(ctx, pool, clearOlderThan, logger)
			}
		}
	}()

	return nil
}

func cleanEvents(ctx context.Context, pool *pgxpool.Pool, clearOlderThan int, logger *libpack_logger.Logger) {
	var errors []error
	var failedQueries []string

	// Format interval parameter for PostgreSQL
	interval := fmt.Sprintf("%d days", clearOlderThan)

	for _, query := range delQueries {
		// Use parameterized query with bound parameter to prevent SQL injection
		_, err := pool.Exec(ctx, query, interval)
		if err != nil {
			errors = append(errors, err)
			failedQueries = append(failedQueries, query)
		} else {
			logger.Debug(&libpack_logger.LogMessage{
				Message: "Successfully executed query",
				Pairs:   map[string]any{"query": query, "interval": interval},
			})
		}
	}

	if len(errors) > 0 {
		var errMsgs []string
		for _, err := range errors {
			errMsgs = append(errMsgs, err.Error())
		}
		logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to execute some queries",
			Pairs: map[string]any{
				"failed_queries": failedQueries,
				"errors":         errMsgs,
			},
		})
	}
}
