package repository

import (
	"context"
	"database/sql"
	"fmt"
	"server/internal/model"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLImageRepo struct {
	db *sql.DB
}

// Khởi tạo repo với db đã kết nối
func NewMySQLImageRepo(db *sql.DB) *MySQLImageRepo {
	return &MySQLImageRepo{db: db}
}
func (r *MySQLImageRepo) Close() error {
	fmt.Println(">>> Closing MySQL connection")
	return r.db.Close()
}

func (r *MySQLImageRepo) Create(ctx context.Context, s *model.Image) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO stories (id, title, author, cover_url, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		s.ID, s.Title, "s.Author", "s.CoverURL", s.CreatedAt,
	)
	return err
}

func (r *MySQLImageRepo) GetByID(ctx context.Context, id string) (*model.Image, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, title, author, cover_url, created_at
		 FROM stories WHERE id = ?`, id,
	)

	var s model.Image
	if err := row.Scan(&s.ID, &s.Title, &s.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *MySQLImageRepo) List(ctx context.Context, offset, limit int) ([]model.Image, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, author, cover_url, created_at
		 FROM stories ORDER BY created_at DESC
		 LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Image
	for rows.Next() {
		var s model.Image
		if err := rows.Scan(&s.ID, &s.Title, &s.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}
