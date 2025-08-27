package model

import "time"

type Story struct {
	ID        string
	Title     string
	Author    string
	CoverURL  string
	CreatedAt time.Time
}

type Chapter struct {
	ID        int64
	StoryID   string
	Title     string
	Content   string
	ImageURL  string
	OrderStt  int64
	CreatedAt time.Time
}

type Image struct {
	ID        string
	Title     string
	ChapterId int64
	Referer   string
	ImageURL  string
	OrderStt  int64
	CreatedAt time.Time
}
