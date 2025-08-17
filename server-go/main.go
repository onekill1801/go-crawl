package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Item struct {
	Referer  string `json:"referer"`
	ImageURL string `json:"image_url"`
	Title    string `json:"title"`
	Order    int    `json:"order"`
}
type InnerData struct {
	Referer  string `json:"referer"`
	ImageURL string `json:"image_url"`
	Title    string `json:"title"`
	Order    int    `json:"order"`
}

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis host:port
		Password: "",               // Password nếu có
		DB:       0,                // DB index
	})

	r := gin.Default()

	// Endpoint GET
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Gin!",
		})
	})

	r.GET("/list", func(c *gin.Context) {
		res, err := rdb.XRange(ctx, "images_queue", "0", "+").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// count := 0
		var items []InnerData
		for _, msg := range res {
			// Lấy field "data" (string JSON)
			dataStr := getString(msg.Values["data"])

			// Parse chuỗi JSON vào struct
			var parsed InnerData
			if err := json.Unmarshal([]byte(dataStr), &parsed); err != nil {
				// Nếu lỗi parse thì bỏ qua entry
				continue
			}
			items = append(items, parsed)
			// count++
			// if count >= 10 {
			// 	break
			// }
		}

		items = append(items, InnerData{
			Referer:  "https://ac.qq.com/",
			ImageURL: "https://manhua.acimg.cn/manhua_detail/0/27_15_39_67ec41e530fbf382fef34251ee200709_1735285149604.jpg/0",
			Title:    "Thêm tay",
			Order:    999,
		})

		c.JSON(http.StatusOK, items)
	})

	// Endpoint GET với param
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{
			"user": name,
		})
	})

	r.GET("/novel/list", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{
			"user": name,
		})
	})

	// Endpoint POST
	r.POST("/data", func(c *gin.Context) {
		var json struct {
			Value string `json:"value" binding:"required"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"received": json.Value,
		})
	})

	r.Run(":8080") // Chạy server ở cổng 8080
}

func getString(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}
