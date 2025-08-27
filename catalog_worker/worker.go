package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"catalog/parser"
	"catalog/queue"

	"golang.org/x/net/html"
)

type DomainJob struct {
	Domain string `json:"domain"`
}

type SeriesJob struct {
	SeriesURL string `json:"series_url"`
	StoryID   string `json:"story_id"`
}

type ChapterJob struct {
	ChapterURL string `json:"chapter_url"`
	Title      string `json:"title"`
	SeriesID   string `json:"series_id"`
	OrderStt   int32  `json:"order_stt"`
	StoryID    string `json:"story_id"`
	ImageUrl   string `json:"image_url"`
}

type ImagesJob struct {
	Referer   string `json:"referer"`
	ImageURL  string `json:"image_url"`
	Title     string `json:"title"`
	OrderStt  int    `json:"order"`
	StoryID   string `json:"story_id"`
	ChapterID int64  `json:"chapter_id"`
}

// func fetchHTML(url string) (*html.Node, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	buf, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return html.Parse(bytes.NewReader(buf))
// }

func fetchHTML(url string) (*html.Node, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// fmt.Printf("resp.StatusCode %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Ghi ra file để kiểm tra nếu cần
	// os.WriteFile("last_response.html", buf, 0644)

	return html.Parse(bytes.NewReader(buf))
}

func ProcessDomainJob(job DomainJob, q *queue.RedisQueue) error {
	parser := parser.GetParserForDomain(job.Domain)
	if parser == nil {
		return fmt.Errorf("no parser for domain: %s", job.Domain)
	}

	rootURL := "https://" + job.Domain
	doc, err := fetchHTML(rootURL)
	if err != nil {
		return err
	}

	seriesList, err := parser.GetSeriesList(doc)
	if err != nil {
		return err
	}

	for _, series := range seriesList {
		id, _ := extractTitleNo(series.LinkNovel)
		q.Push("series_queue", SeriesJob{
			SeriesURL: series.LinkNovel,
			StoryID:   strconv.FormatInt(id, 10),
		})
	}

	return nil
}

func extractTitleNo(url string) (int64, error) {
	re := regexp.MustCompile(`title_no=(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return 0, fmt.Errorf("title_no not found in url: %s", url)
	}
	var num int64
	fmt.Sscanf(matches[1], "%d", &num)
	return num, nil
}

var index = 1

func ProcessSeriesJob(job SeriesJob, q *queue.RedisQueue) (int, error) {
	domain := "www.webtoons.com/en/"
	parser := parser.GetParserForDomain(domain)
	if parser == nil {
		return -1, fmt.Errorf("no parser for domain: %s", domain)
	}

	doc, err := fetchHTML(job.SeriesURL)

	// err = logTextToFile(doc, strconv.Itoa(index)+"_page_log.html")
	// index += 1
	// if err != nil {
	// 	log.Println("error logging HTML:", err)
	// }

	if err != nil {
		return -1, err
	}
	chapters, maxPage, err := parser.GetChapters(doc)
	if err != nil {
		return -1, err
	}

	for _, chap := range chapters {
		orderStt, _ := extractChapterNumber(chap.Title)
		q.Push("chapter_queue", ChapterJob{
			ChapterURL: chap.URL,
			Title:      chap.Title,
			SeriesID:   generateID(job.SeriesURL),
			OrderStt:   orderStt,
		})
	}

	return maxPage, nil
}

func extractChapterNumber(title string) (int32, error) {
	re := regexp.MustCompile(`(?i)^Ep\.?\s*(\d+)`)
	matches := re.FindStringSubmatch(title)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no chapter number found in: %s", title)
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	return int32(num), nil
}

func extractDomain(url string) string {
	// đơn giản hóa
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	return strings.Split(url, "/")[0]
}

func generateID(url string) string {
	// TODO: hash hoặc slugify
	return strings.ReplaceAll(url, "/", "_")
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	text := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return text
}

func logTextToFile(doc *html.Node, filename string) error {
	// Tạo file hoặc ghi đè
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Trích xuất nội dung text từ cây DOM
	textContent := extractText(doc)

	// Ghi nội dung vào file
	_, err = file.WriteString(textContent)
	if err != nil {
		return err
	}

	return nil
}

func ProcessImagesJob(job SeriesJob, q *queue.RedisQueue) error {
	domain := "www.webtoons.com/en/"
	parser := parser.GetParserForDomain(domain)
	if parser == nil {
		return fmt.Errorf("no parser for domain: %s", domain)
	}

	doc, err := fetchHTML(job.SeriesURL)

	if err != nil {
		return err
	}
	images, err := parser.GetListImages(doc)
	if err != nil {
		return err
	}

	for _, img := range images {
		q.Push("images_queue", ImagesJob{
			ImageURL:  img.URL,
			Title:     img.Title,
			OrderStt:  img.Order,
			Referer:   img.Referer,
			StoryID:   "1",
			ChapterID: 554,
		})
	}

	return nil
}
