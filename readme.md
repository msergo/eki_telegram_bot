# 🇪🇪 Estonian Dictionary Telegram Bot
[![Build status](https://dl.circleci.com/status-badge/img/gh/msergo/eki_telegram_bot/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/msergo/eki_telegram_bot/tree/master)

## About

Telegram bot [@eki_ee_bot](https://t.me/eki_ee_bot) for querying translations from official [Estonian dictionaries](http://eki.ee/). Since there is no official API available, the bot parses HTML pages to fetch translations. Currently, it supports Estonian-Russian and Russian-Estonian translations, identifying the direction based on the charset. The bot does not have a permanent storage solution but utilizes Redis to cache already fetched articles. Messages are returned with inline keyboard for quick switching. 

## Motivation
The motivation behind this project was the need for a convenient tool to obtain translations quickly. The user interface of the eki.ee website is far away from being a user-friendly. As Telegram is my daily messenger, choosing it as the platform was an obvious decision. By that time, I had barely any knowledge of Go, so it was a good chance to learn by developing. Surprisingly, the bot works quite well and requires almost no maintenance.

![Screen](./screen.gif)

## Local run
Submit BOT_TOKEN and WEBHOOK_ADDRESS to env in the docker-compose.yaml and run it with `docker-compose up`. 
Now you're able to emulate webhook requests with `POST /localhost:8083`

## Cache Engine Selection

The bot supports two interchangeable cache backends, selected at runtime via the `CACHE_ENGINE` environment variable:

| Value               | Backend | Characteristics                                            |
| ------------------- | ------- | ---------------------------------------------------------- |
| `redis` *(default)* | Redis 7 | Fast in-memory cache; ephemeral (data lost on restart)     |
| `sqlite`            | SQLite  | File-backed persistent cache; no external service required |

### Quickstart: Redis (default)

```bash
docker compose --profile redis-cache up
```

Set `REDIS_HOST`, `BOT_TOKEN`, and `WEBHOOK_ADDRESS` in `docker-compose.yaml` before starting.

### Quickstart: SQLite

```bash
docker compose --profile sqlite-cache up
```

No Redis service is required. The cache database is stored in a Docker volume (`sqlite_data`) and survives container restarts.

### Environment Variables

| Variable          | Required   | Default | Description                               |
| ----------------- | ---------- | ------- | ----------------------------------------- |
| `CACHE_ENGINE`    | No         | `redis` | Cache backend (`redis` or `sqlite`)       |
| `REDIS_HOST`      | Redis only | —       | Hostname of the Redis service             |
| `REDIS_PASS`      | No         | —       | Redis password                            |
| `BOT_TOKEN`       | Yes        | —       | Telegram Bot API token                    |
| `WEBHOOK_ADDRESS` | Yes        | —       | Public HTTPS URL for the Telegram webhook |
| `PORT`            | No         | `8083`  | HTTP port the bot listens on              |

### Switching Engines

Because each engine starts with an empty cache (Redis is ephemeral; SQLite is a fresh file on first run), switching between engines simply requires changing `CACHE_ENGINE`. There is no data migration needed — the cache is a read-through layer and will be repopulated on demand from the EKI website.
