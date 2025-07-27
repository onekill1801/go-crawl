package parser

import "golang.org/x/net/html"

type SeriesURL struct {
	URL  string
	Name string
}

type Series struct {
	LinkNovel string
	Title     string
	Genre     string
	Avatar    string
}

type Chapter struct {
	URL   string
	Title string
}

type ImagesChapter struct {
	NumberChap string
	URL        string
	Title      string
	Referer    string
	Order      int
}

type Parser interface {
	DomainMatch(domain string) bool
	GetSeriesList(doc *html.Node) ([]Series, error)
	GetChapters(doc *html.Node) ([]Chapter, error)
	GetListImages(doc *html.Node) ([]ImagesChapter, error)
}
