package main

import "server/internal/app"

func main() {
	r, repo := app.Setup()
	defer repo.Close() // ✅ Close khi app dừng
	r.Run(":8088")
}
