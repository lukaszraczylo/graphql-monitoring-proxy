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

var delQueries = [...]string{
	"DELETE FROM hdb_catalog.event_invocation_logs WHERE created_at < NOW() - interval '%d days';",
	"DELETE FROM hdb_catalog.event_log WHERE created_at < NOW() - interval '%d days';",
	"DELETE FROM hdb_catalog.hdb_action_log WHERE created_at < NOW() - INTERVAL '%d days';",
	"DELETE FROM hdb_catalog.hdb_cron_event_invocation_logs WHERE created_at < NOW() - INTERVAL '%d days';",
	"DELETE FROM hdb_catalog.hdb_scheduled_event_invocation_logs WHERE created_at < NOW() - INTERVAL '%d days';",
}

func enableHasuraEventCleaner() {
	cfgMutex.RLock()
	if !cfg.HasuraEventCleaner.Enable {
		cfgMutex.RUnlock()
		return
	}

	eventMetadataDb := cfg.HasuraEventCleaner.EventMetadataDb
	if eventMetadataDb == "" {
		logger := cfg.Logger
		cfgMutex.RUnlock()
		
		logger.Warning(&libpack_logger.LogMessage{
			Message: "Event metadata db URL not specified, event cleaner not active",
		})
		return
	}
	
	clearOlderThan := cfg.HasuraEventCleaner.ClearOlderThan
	logger := cfg.Logger
	cfgMutex.RUnlock()

	logger.Info(&libpack_logger.LogMessage{
		Message: "Event cleaner enabled",
		Pairs:   map[string]interface{}{"interval_in_days": clearOlderThan},
	})

	go func(dbURL string, clearOlderThan int, logger *libpack_logger.Logger) {
		pool, err := pgxpool.New(context.Background(), dbURL)
		if err != nil {
			logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to create connection pool",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			return
		}
		defer pool.Close()

		time.Sleep(initialDelay)

		logger.Info(&libpack_logger.LogMessage{
			Message: "Initial cleanup of old events",
		})
		cleanEvents(pool, clearOlderThan, logger)

		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			logger.Info(&libpack_logger.LogMessage{
				Message: "Cleaning up old events",
			})
			cleanEvents(pool, clearOlderThan, logger)
		}
	}(eventMetadataDb, clearOlderThan, logger)
}

func cleanEvents(pool *pgxpool.Pool, clearOlderThan int, logger *libpack_logger.Logger) {
	ctx := context.Background()
	var errors []error
	var failedQueries []string

	for _, query := range delQueries {
		_, err := pool.Exec(ctx, fmt.Sprintf(query, clearOlderThan))
		if err != nil {
			errors = append(errors, err)
			failedQueries = append(failedQueries, query)
		} else {
			logger.Debug(&libpack_logger.LogMessage{
				Message: "Successfully executed query",
				Pairs:   map[string]interface{}{"query": query},
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
			Pairs: map[string]interface{}{
				"failed_queries": failedQueries,
				"errors":         errMsgs,
			},
		})
	}
}
