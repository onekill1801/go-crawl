package repository

import (
	"context"
	"database/sql"
	"fmt"
	"server/internal/model"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLStoryRepo struct {
	db *sql.DB
}

// Khởi tạo repo với db đã kết nối
func NewMySQLStoryRepo() (*MySQLStoryRepo, error) {
	dsn := "root:your_root_password@tcp(192.168.1.6:5306)/story?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open mysql: %w", err)
	}

	// Kiểm tra kết nối thực sự
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping mysql: %w", err)
	}

	return &MySQLStoryRepo{db: db}, nil
}

func (r *MySQLStoryRepo) Close() error {
	fmt.Println(">>> Closing MySQL connection")
	return r.db.Close()
}

func (r *MySQLStoryRepo) Create(ctx context.Context, s *model.Story) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO stories (id, title, author, cover_url, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		s.ID, s.Title, s.Author, s.CoverURL, s.CreatedAt,
	)
	return err
}

func (r *MySQLStoryRepo) GetByID(ctx context.Context, id string) (*model.Story, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, title, author, cover_url, created_at
		 FROM stories WHERE id = ?`, id,
	)

	var s model.Story
	if err := row.Scan(&s.ID, &s.Title, &s.Author, &s.CoverURL, &s.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *MySQLStoryRepo) List(ctx context.Context, offset, limit int) ([]model.Story, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, author, cover_url, created_at
		 FROM stories ORDER BY created_at DESC
		 LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Story
	for rows.Next() {
		var s model.Story
		if err := rows.Scan(&s.ID, &s.Title, &s.Author, &s.CoverURL, &s.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}
