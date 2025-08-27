package parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type TruyenDepParser struct{}

func (p *TruyenDepParser) DomainMatch(domain string) bool {
	return strings.Contains(domain, "www.webtoons.com/en/")
}

func (p *TruyenDepParser) GetSeriesList(doc *html.Node) ([]Series, error) {
	// TODO get ListTreding
	var results []Series

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var series Series
			var isTarget bool

			// Kiểm tra class
			for _, attr := range n.Attr {
				if attr.Key == "class" && (strings.Contains(attr.Val, "link _trending_title_a") ||
					strings.Contains(attr.Val, "link _genre_title_a") ||
					strings.Contains(attr.Val, "link _canvas_title_a") ||
					strings.Contains(attr.Val, "link _daily_title_a") ||
					strings.Contains(attr.Val, "link _new_title_a")) {
					isTarget = true
				}
				if attr.Key == "href" {
					series.LinkNovel = attr.Val
				}
			}

			if isTarget {
				// Duyệt toàn bộ con
				var extract func(*html.Node)
				extract = func(child *html.Node) {
					if child.Type == html.ElementNode {
						// Avatar
						if child.Data == "img" {
							for _, attr := range child.Attr {
								if attr.Key == "src" {
									series.Avatar = attr.Val
								}
							}
						}
						// Title
						if child.Data == "strong" {
							for _, attr := range child.Attr {
								if attr.Key == "class" && attr.Val == "title" {
									if child.FirstChild != nil {
										series.Title = strings.TrimSpace(child.FirstChild.Data)
									}
								}
							}
						}
						// Genre
						if child.Data == "div" {
							for _, attr := range child.Attr {
								if attr.Key == "class" && attr.Val == "genre" {
									if child.FirstChild != nil {
										series.Genre = strings.TrimSpace(child.FirstChild.Data)
									}
								}
							}
						}
					}
					// tiếp tục duyệt con
					for c := child.FirstChild; c != nil; c = c.NextSibling {
						extract(c)
					}
				}
				extract(n)

				// Thêm nếu có đủ dữ liệu
				if series.LinkNovel != "" && series.Title != "" {
					results = append(results, series)
				}
			}
		}

		// duyệt tiếp các node khác
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return results, nil
}

func (p *TruyenDepParser) GetChapters(doc *html.Node) ([]Chapter, int, error) {

	maxPage := extractMaxPage(doc)
	fmt.Println("Total pages:", maxPage)

	var results []Chapter
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "detail_lst" {
					// Bắt đầu tìm ul#_listUl bên trong div này
					findUl(n, &results)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return results, maxPage, nil
}

func findUl(n *html.Node, links *[]Chapter) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "ul" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "_listUl" {
					// Đã tìm được ul#_listUl, xử lý tiếp các li bên trong
					extractChaptersFromUl(n, links)
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
}

func extractChaptersFromUl(ul *html.Node, chapters *[]Chapter) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			var chapter Chapter

			// Tìm thẻ <a> trong <li>
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "a" {
					// Lấy href
					for _, attr := range c.Attr {
						if attr.Key == "href" {
							chapter.URL = attr.Val
						}
					}

					// Tìm <span class="subj"> trong <a>
					for aChild := c.FirstChild; aChild != nil; aChild = aChild.NextSibling {
						if aChild.Type == html.ElementNode && aChild.Data == "span" {
							for _, attr := range aChild.Attr {
								if attr.Key == "class" && attr.Val == "subj" {
									chapter.Title = extractFirstSpanText(aChild)
								}
							}
						}
					}
				}
			}

			if chapter.URL != "" {
				*chapters = append(*chapters, chapter)
			}
		}

		// Tiếp tục duyệt các li khác
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(ul)
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractText(c) + " "
	}
	return strings.TrimSpace(result)
}

func extractFirstSpanText(subj *html.Node) string {
	for c := subj.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "span" {
			// Lấy text đầu tiên trong <span>
			return extractText(c)
		}
	}
	return ""
}

func extractMaxPage(doc *html.Node) int {
	maxPage := 0

	var findPaginate func(*html.Node)
	findPaginate = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "paginate" {
					// Đã tìm được <div class="paginate">
					extractPageNumbers(n, &maxPage)
					return
				}
			}
		}

		// Duyệt tiếp
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findPaginate(c)
		}
	}

	findPaginate(doc)
	return maxPage
}

func extractPageNumbers(paginateDiv *html.Node, maxPage *int) {
	re := regexp.MustCompile(`page=(\d+)`)
	var lastHref string
	hasNext := false

	for c := paginateDiv.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			var href string
			var class string

			for _, attr := range c.Attr {
				if attr.Key == "href" {
					href = attr.Val
				}
				if attr.Key == "class" {
					class = attr.Val
				}
			}

			if class == "pg_next" {
				hasNext = true
				break
			}

			if href != "" && href != "#" {
				lastHref = href
			}
		}
	}

	if hasNext {
		*maxPage = 0
		return
	}

	if lastHref != "" {
		if matches := re.FindStringSubmatch(lastHref); len(matches) == 2 {
			if page, err := strconv.Atoi(matches[1]); err == nil {
				*maxPage = page
				return
			}
		}
	}

	// nếu chỉ có 1 page hoặc href = '#' thì mặc định là 1
	*maxPage = 1
}

func extractPageNumberFromHref(href string) int {
	u, err := url.Parse(href)
	if err != nil {
		return 0
	}

	pageStr := u.Query().Get("page")
	if pageStr == "" {
		return 0
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0
	}

	return page
}

func (p *TruyenDepParser) GetListImages(doc *html.Node) ([]ImagesChapter, error) {
	var results []ImagesChapter
	var order int = 1

	var findImages func(*html.Node)
	findImages = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			var hasTargetClass, hasTargetID bool
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "viewer_img _img_viewer_area") {
					hasTargetClass = true
				}
				if attr.Key == "id" && attr.Val == "_imageList" {
					hasTargetID = true
				}
			}

			if hasTargetClass && hasTargetID {
				// Tìm thẻ <img> bên trong div này
				var extractImages func(*html.Node)
				extractImages = func(child *html.Node) {
					if child.Type == html.ElementNode && child.Data == "img" {
						for _, attr := range child.Attr {
							if attr.Key == "data-url" && strings.TrimSpace(attr.Val) != "" {
								results = append(results, ImagesChapter{
									NumberChap: "numberChap",
									URL:        attr.Val,
									Title:      "title",
									Order:      order,
									Referer:    "https://www.webtoons.com",
									StoryID:    "storyID",
									ChapterID:  1,
								})
								order++
								break
							}
						}
					}
					for c := child.FirstChild; c != nil; c = c.NextSibling {
						extractImages(c)
					}
				}
				extractImages(n)
			}
		}

		// tiếp tục duyệt toàn bộ cây DOM
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findImages(c)
		}
	}

	findImages(doc)

	return results, nil
}
