package main

import (
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

func (r RedisWorker) Ping() (string, error) {
	return r.client.Ping().Result()
}

func (r RedisWorker) StoreArticles(key string, articles []string) error {
	for i := len(articles) - 1; i >= 0; i-- {
		if err := r.client.LPush(key, articles[i]).Err(); err != nil {
			return err
		}
	}
	return nil
}

// Returns an empty slice when the key does not exist.
func (r RedisWorker) GetAllArticles(key string) ([]string, error) {
	res := r.client.LRange(key, 0, -1)
	return res.Val(), res.Err()
}

// Returns an empty string when the index is out of range.
func (r RedisWorker) GetArticleByIndex(key string, index int64) (string, error) {
	res := r.client.LIndex(key, index)
	// Redis returns redis.Nil when the index does not exist; treat that as empty.
	if res.Err() == redis.Nil {
		return "", nil
	}
	return res.Val(), res.Err()
}

func (r RedisWorker) GetArticlesLen(key string) (int, error) {
	res := r.client.LLen(key)
	return int(res.Val()), res.Err()
}

func (r RedisWorker) Close() error {
	return r.client.Close()
}
