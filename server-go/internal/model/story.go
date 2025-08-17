package model

import "time"

type Story struct {
	ID        string
	Title     string
	Author    string
	CoverURL  string
	CreatedAt time.Time
}
