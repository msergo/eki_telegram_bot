package main

import (
	env "github.com/Netflix/go-env"
)

type Environment struct {
	RedisHost      string `env:"REDIS_HOST"`
	BotToken       string `env:"BOT_TOKEN"`
	WebhookAddress string `env:"WEBHOOK_ADDRESS"`
	AppPort        string `env:"PORT"`
	SentryDsn      string `env:"SENTRY_DSN"`
	Env            string `env:"ENV"`
	RedisPass      string `env:"REDIS_PASS"`
	// CacheEngine selects the active cache backend. Valid values: "sqlite" (default), "redis".
	CacheEngine string `env:"CACHE_ENGINE,default=sqlite"`
	// SqliteDbPath specifies the path to the SQLite database file. Default: "./cache.db".
	SqliteDbPath string `env:"SQLITE_DB_PATH,default=./cache.db"`

	Extras env.EnvSet
}
