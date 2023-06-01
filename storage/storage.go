package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type Quote struct {
	Text       string
	CategoryId int
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
	q := `CREATE TABLE IF NOT EXISTS quotes (quote_text TEXT, category_id int, username TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}

func (s *Storage) Add(ctx context.Context, quote *Quote) error {
	q := `INSERT INTO quotes (quote_text, category_id, username) VALUES (?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, quote.Text, quote.CategoryId, quote.Username); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, categoryId int, username string) (*Quote, error) {
	q := `SELECT quote_text FROM quotes WHERE username = ? AND categoryId = ? ORDER BY RANDOM() LIMIT 1`

	var quote string
	err := s.db.QueryRowContext(ctx, q, categoryId, username).Scan(&quote)
	if err == sql.ErrNoRows {
		return nil, errors.New("no saved quotes")
	}

	if err != nil {
		return nil, fmt.Errorf("can't pick random quote: %w", err)
	}
	return &Quote{Text: quote, Username: username}, nil
}

func (s *Storage) IfExists(ctx context.Context, quote *Quote) (bool, error) {
	q := `SELECT COUNT(*) FROM quotes WHERE username = ? AND quote_text = ? AND category_id = ?`

	var count int
	err := s.db.QueryRowContext(ctx, q, quote.Username, quote.Text, quote.CategoryId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("can't check if the quote exists: %w", err)
	}

	return count > 0, err
}
