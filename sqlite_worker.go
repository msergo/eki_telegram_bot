package main

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteWorker struct {
	db *sql.DB
}

func InitSQLiteWorker(dbPath string) (SQLiteWorker, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return SQLiteWorker{}, err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return SQLiteWorker{}, err
	}

	err = runMigrations(db)
	if err != nil {
		return SQLiteWorker{}, err
	}

	return SQLiteWorker{db: db}, nil
}

// TODO add proper migrations later
func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS articles (
		search_key TEXT PRIMARY KEY,
		articles_json TEXT
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_articles_search_key ON articles(search_key)`)
	return err
}

// StoreArticles serialises articles as JSON and upserts the row for key.
func (s SQLiteWorker) StoreArticles(key string, articles []string) error {
	data, err := json.Marshal(articles)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT OR REPLACE INTO articles (search_key, articles_json) VALUES (?, ?)`, key, string(data))
	return err
}

// GetAllArticles returns an empty slice (not an error) when key does not exist.
func (s SQLiteWorker) GetAllArticles(key string) ([]string, error) {
	var data string
	err := s.db.QueryRow(`SELECT articles_json FROM articles WHERE search_key = ?`, key).Scan(&data)
	if err == sql.ErrNoRows {
		return []string{}, nil
	}
	if err != nil {
		return []string{}, err
	}
	var articles []string
	if err := json.Unmarshal([]byte(data), &articles); err != nil {
		return []string{}, err
	}
	return articles, nil
}

// GetArticleByIndex returns empty string (not an error) when index is out of range.
func (s SQLiteWorker) GetArticleByIndex(key string, index int64) (string, error) {
	articles, err := s.GetAllArticles(key)
	if err != nil {
		return "", err
	}
	if int(index) >= len(articles) || index < 0 {
		return "", nil
	}
	return articles[index], nil
}

func (s SQLiteWorker) GetArticlesLen(key string) (int, error) {
	articles, err := s.GetAllArticles(key)
	return len(articles), err
}

func (s SQLiteWorker) Close() error {
	return s.db.Close()
}
