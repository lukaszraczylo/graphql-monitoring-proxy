package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

func enableHasuraEventCleaner() {
	if cfg.HasuraEventCleaner.Enable {
		if cfg.HasuraEventCleaner.EventMetadataDb == "" {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Event metadata db URL not specified, event cleaner not active",
				Pairs:   nil,
			})
			return
		}

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Event cleaner enabled",
			Pairs:   map[string]interface{}{"interval_in_days": cfg.HasuraEventCleaner.ClearOlderThan},
		})

		time.Sleep(60 * time.Second) // wait for everything to start and settle down
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Initial cleanup of old events",
			Pairs:   nil,
		})
		cleanEvents()

		for {
			select {
			case <-ticker.C:
				cfg.Logger.Info(&libpack_logger.LogMessage{
					Message: "Cleaning up old events",
					Pairs:   nil,
				})
				cleanEvents()
			}
		}
	}
}

func cleanEvents() {
	conn, err := pgx.Connect(context.Background(), cfg.HasuraEventCleaner.EventMetadataDb)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to connect to event metadata db",
			Pairs:   map[string]interface{}{"error": err},
		})
		return
	}
	defer conn.Close(context.Background())

	delQueries := []string{
		fmt.Sprintf("DELETE FROM hdb_catalog.event_invocation_logs WHERE created_at < now() - interval '%d days';", cfg.HasuraEventCleaner.ClearOlderThan),
		fmt.Sprintf("DELETE FROM hdb_catalog.event_log WHERE created_at < now() - interval '%d days';", cfg.HasuraEventCleaner.ClearOlderThan),
		fmt.Sprintf("DELETE FROM hdb_catalog.hdb_action_log WHERE created_at < NOW() - INTERVAL '%d days';", cfg.HasuraEventCleaner.ClearOlderThan),
		fmt.Sprintf("DELETE FROM hdb_catalog.hdb_cron_event_invocation_logs WHERE created_at < NOW() - INTERVAL '%d days';", cfg.HasuraEventCleaner.ClearOlderThan),
		fmt.Sprintf("DELETE FROM hdb_catalog.hdb_scheduled_event_invocation_logs WHERE created_at < NOW() - INTERVAL '%d days';", cfg.HasuraEventCleaner.ClearOlderThan),
	}

	for _, query := range delQueries {
		_, err := conn.Exec(context.Background(), query)
		if err != nil {
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Failed to execute query",
				Pairs:   map[string]interface{}{"query": query, "error": err},
			})
		}
	}
	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Old events cleaned up",
		Pairs:   nil,
	})
}
