package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

type Quote struct {
	Text       string
	Author     string
	CategoryId int
	GenderId   int
	Username   string
}

func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS quotes (quote_text TEXT, author TEXT, category_id int, genderId int, username TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}

func (s *Storage) Add(ctx context.Context, quote *Quote) error {
	q := `INSERT INTO quotes (quote_text, author, category_id, gender_id, username) VALUES (?, ?, ?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, quote.Text, quote.Author, quote.CategoryId, quote.GenderId, quote.Username); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, categoryId, genderId int, username string) (*Quote, error) {
	q := `SELECT quote_text FROM quotes WHERE username = ? AND categoryId = ? AND genderId = ? ORDER BY RANDOM() LIMIT 1`

	var quote string
	err := s.db.QueryRowContext(ctx, q, categoryId, genderId, username).Scan(&quote)
	if err == sql.ErrNoRows {
		return nil, errors.New("no saved pages")
	}

	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &Quote{Text: quote, Username: username}, nil
}
