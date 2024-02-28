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

	Extras env.EnvSet
}
