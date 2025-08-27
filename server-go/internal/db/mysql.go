// internal/db/mysql.go
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL() (*sql.DB, error) {
	dsn := "root:your_root_password@tcp(192.168.1.6:5306)/story?parseTime=true"
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
