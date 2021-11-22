package page

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Page interface {
	GetTitle() string
	GetLinks() []string
}

type page struct {
	doc *goquery.Document
}

func NewPage(raw io.Reader) (Page, error) {
	doc, err := goquery.NewDocumentFromReader(raw)
	if err != nil {
		return nil, err
	}
	return &page{doc: doc}, nil
}

func (p *page) GetTitle() string {
	return p.doc.Find("title").First().Text()
}

func (p *page) GetLinks() []string {
	var urls []string
	p.doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			 if strings.HasPrefix(url, "/") {
				url = fmt.Sprintf( "https://www.telegram.com%s", url)
				urls = append(urls, url)
				 return
			}
			if strings.HasPrefix(url, "mailto") {
				return
			}
			if strings.HasPrefix(url, "#") {
				return
			}
				urls = append(urls, url)
		}
	})
	return urls
}
