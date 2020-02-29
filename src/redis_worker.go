package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"strconv"
)

const pubSubChan = "searches"

type RedisWorker struct {
	client *redis.Client
}

func InitRedisWorker() RedisWorker {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", environment.RedisHost, "6379"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return RedisWorker{client: client}
}

func (r RedisWorker) Ping() (response string, error error) {
	r.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", environment.RedisHost, "6379"),
		Password: "", // no password set
		DB:       0,  // use default DB
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
	r.client.Expire(key, time.Hour*48)
}

func (r RedisWorker) GetAllArticles(key string) []string {
	return r.client.LRange(key, 0, -1).Val()
}

func (r RedisWorker) GetArticleByIndex(key string, indexStr string) string {
	index, _ := strconv.ParseInt(indexStr, 10, 64)
	return r.client.LIndex(key, index).Val()
}

func (r RedisWorker) GetArticlesLenByKeyword(key string) int {
	len64 := r.client.LLen(key).Val()
	return int(len64)
}

func (r RedisWorker) pushToChannel(value string) error {
	return r.client.Publish(pubSubChan, value).Err()
}
