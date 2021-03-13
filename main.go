package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

const START_URL = "https://en.wikipedia.org"

func main() {
	res, err := http.Get(START_URL)

	if err != nil {
		panic(err)
	}

	links := findLinks(res.Body).AddDomain(START_URL).Filter().Clean()

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
