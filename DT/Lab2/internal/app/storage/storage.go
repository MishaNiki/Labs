package storage

import (
	"database/sql"
	"fmt"

	"github.com/MishaNiki/Labs/DT/Lab2/internal/app/spamdetection"
	_ "github.com/lib/pq" // driver postgresql
)

// Storage ...
type Storage struct {
	db *sql.DB
}

// Config ...
type Config struct {
	User     string `json:"pg-user"`
	Password string `json:"pg-password"`
	DBName   string `json:"pg-dbname"`
	SSLMode  string `json:"pg-sslmode"`
}

// New ...
func New() *Storage {
	return &Storage{}
}

// Open ...
func (s *Storage) Open(config *Config) error {

	dbURL := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}

// Close ...
func (s *Storage) Close() {
	s.db.Close()
}

// LoadStatistics ...
func (s *Storage) LoadStatistics(dict *spamdetection.Dictionary) {
	data := dict.String()
	s.db.QueryRow("SELECT \"DTM\".\"AddStatistics\"($1)", data)
}

// DownloadStatistics ...
func (s *Storage) DownloadStatistics() (*spamdetection.Dictionary, error) {

	rows, err := s.db.Query("SELECT * FROM \"DTM\".\"Lab2Stat\"")
	if err != nil {
		return nil, err
	}

	dict := spamdetection.NewDictionary()

	for rows.Next() {
		word := spamdetection.Word{}
		err := rows.Scan(&word.Word, &word.AmountOK, &word.AmountSpam)
		if err != nil {
			return nil, err
		}
		dict.AddWord(word)
	}

	return dict, nil
}
