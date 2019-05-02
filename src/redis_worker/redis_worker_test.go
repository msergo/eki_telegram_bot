package redis_worker

import (
	"testing"
)

func TestRedisWorker_Ping(t *testing.T) {
	w := InitRedisWorker()
	defer w.client.Close()
	pong, _ := w.Ping()
	if (pong != "PONG") {
		t.Fail()
	}
}

func TestRedisWorker_StoreArticles(t *testing.T) {
	w := InitRedisWorker()
	defer w.client.Close()
	data := []string{"Japan", "Australia", "Germany"}
	w.StoreArticlesSet(123, data)
	australiaElem := w.GetArticleByIndex(123, 1)
	if australiaElem != "Australia" {
		t.Fail()
	}
}
