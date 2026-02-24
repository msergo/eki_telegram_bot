package main

import (
	"os"
	"testing"
)

// createSQLiteTestEngine returns a SQLiteWorker backed by a temporary file.
// The caller is responsible for calling engine.Close() and removing the file.
func createSQLiteTestEngine(t *testing.T) (CacheEngine, func()) {
	t.Helper()
	f, err := os.CreateTemp("", "eki_test_*.db")
	if err != nil {
		t.Fatalf("failed to create temp db file: %v", err)
	}
	f.Close()

	engine, err := InitSQLiteWorker(f.Name())
	if err != nil {
		os.Remove(f.Name())
		t.Fatalf("InitSQLiteWorker: %v", err)
	}
	cleanup := func() {
		engine.Close()
		os.Remove(f.Name())
	}
	return engine, cleanup
}

// engineCases returns the engines available for parameterized testing.
// Redis is included only when REDIS_HOST is set in the environment.
func engineCases(t *testing.T) []struct {
	name   string
	engine CacheEngine
	done   func()
} {
	t.Helper()
	sqEngine, sqCleanup := createSQLiteTestEngine(t)
	cases := []struct {
		name   string
		engine CacheEngine
		done   func()
	}{
		{"sqlite", sqEngine, sqCleanup},
	}

	if host := os.Getenv("REDIS_HOST"); host != "" {
		environment.RedisHost = host
		environment.RedisPass = os.Getenv("REDIS_PASS")
		rw := InitRedisWorker()
		if _, err := rw.Ping(); err == nil {
			cases = append(cases, struct {
				name   string
				engine CacheEngine
				done   func()
			}{"redis", rw, func() { rw.Close() }})
		}
	}

	return cases
}

// ---------------------------------------------------------------------------
// Shared behavioural tests
// ---------------------------------------------------------------------------

func testStoreAndRetrieve(t *testing.T, engine CacheEngine) {
	t.Helper()
	const key = "test_store"
	want := []string{"article1", "article2", "article3"}

	if err := engine.StoreArticles(key, want); err != nil {
		t.Fatalf("StoreArticles: %v", err)
	}

	got, err := engine.GetAllArticles(key)
	if err != nil {
		t.Fatalf("GetAllArticles: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("GetAllArticles: got %d articles, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("GetAllArticles[%d]: got %q, want %q", i, got[i], want[i])
		}
	}
}

func testIndexAccess(t *testing.T, engine CacheEngine) {
	t.Helper()
	const key = "test_index"
	articles := []string{"first", "second", "third"}

	if err := engine.StoreArticles(key, articles); err != nil {
		t.Fatalf("StoreArticles: %v", err)
	}

	for i, want := range articles {
		got, err := engine.GetArticleByIndex(key, int64(i))
		if err != nil {
			t.Fatalf("GetArticleByIndex(%d): %v", i, err)
		}
		if got != want {
			t.Errorf("GetArticleByIndex(%d): got %q, want %q", i, got, want)
		}
	}
}

func testEmptyKey(t *testing.T, engine CacheEngine) {
	t.Helper()
	const key = "nonexistent_key_xyz"

	articles, err := engine.GetAllArticles(key)
	if err != nil {
		t.Fatalf("GetAllArticles on missing key: %v", err)
	}
	if len(articles) != 0 {
		t.Errorf("expected empty slice, got %v", articles)
	}

	n, err := engine.GetArticlesLen(key)
	if err != nil {
		t.Fatalf("GetArticlesLen on missing key: %v", err)
	}
	if n != 0 {
		t.Errorf("expected len 0, got %d", n)
	}

	text, err := engine.GetArticleByIndex(key, 0)
	if err != nil {
		t.Fatalf("GetArticleByIndex on missing key: %v", err)
	}
	if text != "" {
		t.Errorf("expected empty string, got %q", text)
	}
}

func testLargeDataSet(t *testing.T, engine CacheEngine) {
	t.Helper()
	const key = "test_large"
	const n = 50
	articles := make([]string, n)
	for i := range articles {
		articles[i] = "article_content_" + string(rune('A'+i%26))
	}

	if err := engine.StoreArticles(key, articles); err != nil {
		t.Fatalf("StoreArticles: %v", err)
	}

	length, err := engine.GetArticlesLen(key)
	if err != nil {
		t.Fatalf("GetArticlesLen: %v", err)
	}
	if length != n {
		t.Errorf("GetArticlesLen: got %d, want %d", length, n)
	}
}

// ---------------------------------------------------------------------------
// Parameterised runners
// ---------------------------------------------------------------------------

func TestCacheEngine_StoreAndRetrieve(t *testing.T) {
	for _, tc := range engineCases(t) {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			defer tc.done()
			testStoreAndRetrieve(t, tc.engine)
		})
	}
}

func TestCacheEngine_IndexAccess(t *testing.T) {
	for _, tc := range engineCases(t) {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			defer tc.done()
			testIndexAccess(t, tc.engine)
		})
	}
}

func TestCacheEngine_EmptyKey(t *testing.T) {
	for _, tc := range engineCases(t) {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			defer tc.done()
			testEmptyKey(t, tc.engine)
		})
	}
}

func TestCacheEngine_LargeDataSet(t *testing.T) {
	for _, tc := range engineCases(t) {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			defer tc.done()
			testLargeDataSet(t, tc.engine)
		})
	}
}
