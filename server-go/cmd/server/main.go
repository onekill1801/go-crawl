package main

import "server/internal/app"

func main() {
	r := app.Setup()
	r.Run(":8080")
}
