package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func enableHasuraEventCleaner() {
	if cfg.HasuraEventCleaner.Enable {
		if cfg.HasuraEventCleaner.EventMetadataDb == "" {
			cfg.Logger.Warning("Event metadata db URL not specified, event cleaner not active", nil)
			return
		}

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		cfg.Logger.Info("Event cleaner enabled", map[string]interface{}{"interval_in_days": cfg.HasuraEventCleaner.ClearOlderThan})

		go func() {
			for {
				select {
				case <-ticker.C:
					cfg.Logger.Info("Cleaning up old events", nil)
					cleanEvents()
				}
			}
		}()
	}
}

func cleanEvents() {
	conn, err := pgx.Connect(context.Background(), cfg.HasuraEventCleaner.EventMetadataDb)
	if err != nil {
		cfg.Logger.Error("Failed to connect to event metadata db", map[string]interface{}{"error": err})
		return
	}
	defer conn.Close(context.Background())

	delQueries := []string{
		fmt.Sprintf("DELETE FROM hdb_catalog.event_invocation_logs WHERE created_at < now() - interval '%d days';", cfg.HasuraEventCleaner.ClearOlderThan),
		fmt.Sprintf("DELETE FROM hdb_catalog.event_log WHERE created_at < now() - interval '%d days';", cfg.HasuraEventCleaner.ClearOlderThan),
	}

	for _, query := range delQueries {
		_, err := conn.Exec(context.Background(), query)
		if err != nil {
			cfg.Logger.Error("Failed to execute query", map[string]interface{}{"query": query, "error": err})
		}
	}
}
