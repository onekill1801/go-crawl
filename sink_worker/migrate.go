package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql" // <- quan trọng
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func migrate1(dsn string) {
	// dsn := "root:your_root_password@tcp(192.168.1.6:5306)?multiStatements=true"

	m, err := migrate.New(
		"file://./migrations",
		"mysql://"+dsn,
	)
	if err != nil {
		log.Fatal("Cannot create migrate instance:", err)
	}

	// Up migration (tất cả migration chưa chạy)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	// Nếu muốn rollback 1 step
	// if err := m.Steps(-1); err != nil {
	//     log.Fatal(err)
	// }

	log.Println("Migration completed!")
}
