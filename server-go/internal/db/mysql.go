// internal/db/mysql.go
package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL() (*sql.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:crawl_secret@tcp(localhost:3306)/story?parseTime=true"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open mysql: %w", err)
	}

	// kiểm tra kết nối
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping mysql: %w", err)
	}

	// Optional: config pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	return db, nil
}
