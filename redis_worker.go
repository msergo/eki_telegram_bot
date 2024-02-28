package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type RedisWorker struct {
	client *redis.Client
}

func InitRedisWorker() RedisWorker {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", environment.RedisHost, "6379"),
		Password: environment.RedisPass,
		DB:       0,
	})
	return RedisWorker{client: client}
}

func (r RedisWorker) Ping() (response string, error error) {
	r.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", environment.RedisHost, "6379"),
		Password: environment.RedisPass,
		DB:       0,
	})
	pong, err := r.client.Ping().Result()
	return pong, err
}

func (r RedisWorker) StoreArticles(key string, coll []string) error {
	pages, _ := json.Marshal(coll)
	return r.client.LPush(key, pages).Err()
}

func (r RedisWorker) StoreArticlesSet(key string, articles []string) {
	for i := len(articles) - 1; i >= 0; i-- {
		r.client.LPush(key, articles[i]).Err()
	}
}

func (r RedisWorker) GetAllArticles(key string) []string {
	return r.client.LRange(key, 0, -1).Val()
}

func (r RedisWorker) GetArticleByIndex(key string, index int64) string {
	return r.client.LIndex(key, index).Val()
}

func (r RedisWorker) GetArticlesLen(key string) int {
	len64 := r.client.LLen(key).Val()
	return int(len64)
}
