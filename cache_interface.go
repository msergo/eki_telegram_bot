package main

// CacheEngine defines the common interface for all cache storage backends.
// Both Redis and SQLite implementations must satisfy this interface,
// enabling runtime selection via the CACHE_ENGINE environment variable.
type CacheEngine interface {
	// StoreArticles persists a list of article HTML strings under the given key.
	StoreArticles(key string, articles []string) error

	// GetAllArticles retrieves all articles stored under the given key.
	// Returns an empty slice (not an error) when the key does not exist.
	GetAllArticles(key string) ([]string, error)

	// GetArticleByIndex retrieves a single article at the specified position.
	// Returns an empty string (not an error) when the index is out of range.
	GetArticleByIndex(key string, index int64) (string, error)

	// GetArticlesLen returns the number of articles stored under the given key.
	GetArticlesLen(key string) (int, error)

	// Close releases any resources held by the engine (connections, file handles, etc.).
	Close() error
}
