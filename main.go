package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Assumtions:
// - Prepended URL in cleanup step is always the Wikipedia start URL
// - HTTPS used exclusively

func NewLink(url string) Link {
	return Link{
		Url: url,
	}
}

type Link struct {
	Url string
}

func (l Link) Clean() Link {
	// Remove "//" prefix
	l.Url = strings.TrimPrefix(l.Url, "//")

	// Pre-pend website domain name to keys with "/"
	if strings.HasPrefix(l.Url, "/") {
		l.Url = fmt.Sprintf("%s%s", START_URL, l.Url)
	}

	// Pre-pend HTTPS protocol if not exists
	if !strings.HasPrefix(l.Url, "http://") && !strings.HasPrefix(l.Url, "https://") {
		l.Url = fmt.Sprintf("%s%s", "https://", l.Url)
	}

	return l
}

type Links map[Link]bool

func (links Links) Clean() (cleanLinks Links) {
	cleanLinks = make(Links)

	for link := range links {
		cleanLinks[link.Clean()] = true
	}

	return cleanLinks
}

func (links Links) Filter() Links {
	for link := range links {
		if strings.HasPrefix(link.Url, "#") {
			delete(links, link)
		}
	}

	return links
}

func (links Links) String() string {
	linkSlice := make([]string, 0)
	for link := range links {
		linkSlice = append(linkSlice, link.Url)
	}

	return strings.Join(linkSlice, "\n")
}

const START_URL = "https://en.wikipedia.org"

func main() {
	res, err := http.Get(START_URL)

	if err != nil {
		panic(err)
	}

	links := findLinks(res.Body).Filter().Clean()

	fmt.Printf("%v", links)
}

func findLinks(body io.ReadCloser) (links Links) {
	tokenizer := html.NewTokenizer(body)

	links = make(map[Link]bool)

	for {
		if tag := tokenizer.Next(); tag == html.ErrorToken {
			return
		}

		token := tokenizer.Token()

		if token.Data == "a" {
			for _, attr := range token.Attr {

				if attr.Key == "href" {
					link := NewLink(attr.Val)
					links[link] = true
				}
			}
		}
	}
}
