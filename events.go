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

var delQueries = []string{
	"DELETE FROM hdb_catalog.event_invocation_logs WHERE created_at < now() - interval '%d days';",
	"DELETE FROM hdb_catalog.event_log WHERE created_at < now() - interval '%d days';",
	"DELETE FROM hdb_catalog.hdb_action_log WHERE created_at < NOW() - INTERVAL '%d days';",
	"DELETE FROM hdb_catalog.hdb_cron_event_invocation_logs WHERE created_at < NOW() - INTERVAL '%d days';",
	"DELETE FROM hdb_catalog.hdb_scheduled_event_invocation_logs WHERE created_at < NOW() - INTERVAL '%d days';",
}

func enableHasuraEventCleaner() {
	if !cfg.HasuraEventCleaner.Enable {
		return
	}

	if cfg.HasuraEventCleaner.EventMetadataDb == "" {
		cfg.Logger.Warning(&libpack_logger.LogMessage{
			Message: "Event metadata db URL not specified, event cleaner not active",
		})
		return
	}

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Event cleaner enabled",
		Pairs:   map[string]interface{}{"interval_in_days": cfg.HasuraEventCleaner.ClearOlderThan},
	})

	pool, err := pgxpool.New(context.Background(), cfg.HasuraEventCleaner.EventMetadataDb)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to create connection pool",
			Pairs:   map[string]interface{}{"error": err},
		})
		return
	}
	defer pool.Close()

	go func() {
		time.Sleep(initialDelay)

		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Initial cleanup of old events",
		})
		cleanEvents(pool)

		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			cfg.Logger.Info(&libpack_logger.LogMessage{
				Message: "Cleaning up old events",
			})
			cleanEvents(pool)
		}
	}()
}

func cleanEvents(pool *pgxpool.Pool) {
	ctx := context.Background()
	for _, query := range delQueries {
		_, err := pool.Exec(ctx, fmt.Sprintf(query, cfg.HasuraEventCleaner.ClearOlderThan))
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to execute query",
				Pairs:   map[string]interface{}{"query": query, "error": err},
			})
		}
	}
}
