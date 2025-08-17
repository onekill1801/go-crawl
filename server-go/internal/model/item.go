package model

type Item struct {
	Referer  string `json:"referer"`
	ImageURL string `json:"image_url"`
	Title    string `json:"title"`
	Order    int    `json:"order"`
}
