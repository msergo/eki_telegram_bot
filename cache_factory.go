package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// NewCacheEngine returns a CacheEngine implementation based on the engineType string.
// Supported values are "sqlite" (default) and "redis".
// For "sqlite", the database file is opened and migrations are applied.
// For "redis", the factory verifies connectivity via Ping before returning.
func NewCacheEngine(engineType string) (CacheEngine, error) {
	log.WithFields(log.Fields{
		"event_type":   "app_event",
		"cache_engine": engineType,
	}).Info("Initializing cache engine")

	switch engineType {
	case "redis":
		worker := InitRedisWorker()
		_, err := worker.Ping()
		if err != nil {
			return nil, fmt.Errorf("redis ping failed: %w", err)
		}
		log.Info("Cache engine: Redis (ephemeral)")
		return worker, nil

	case "sqlite", "":
		engine, err := InitSQLiteWorker("./cache.db")
		if err != nil {
			return nil, fmt.Errorf("sqlite init failed: %w", err)
		}
		log.Info("Cache engine: SQLite (persistent)")
		return engine, nil

	default:
		return nil, fmt.Errorf("unknown cache engine %q – valid values are \"sqlite\" and \"redis\"", engineType)
	}
}
