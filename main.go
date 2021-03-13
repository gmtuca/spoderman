package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Assumtions:
// - Prepended URL in cleanup step is always the Wikipedia start URL
// - HTTPS used exclusively

const START_URL = "https://en.wikipedia.org"

func main() {
	res, err := http.Get(START_URL)

	if err != nil {
		panic(err)
	}

	urlMap := findLinks(res.Body)
	cleanedUrlMap := cleanLinks(urlMap)

	prettyPrint(cleanedUrlMap)
}

func findLinks(body io.ReadCloser) (urlMap map[string]bool) {
	tokenizer := html.NewTokenizer(body)

	urlMap = make(map[string]bool)

	for {
		if tag := tokenizer.Next(); tag == html.ErrorToken {
			return
		}

		token := tokenizer.Token()

		if token.Data == "a" {
			for _, attr := range token.Attr {
				// formattedUrl := strings.TrimPrefix(attr.Val, "//")

				if attr.Key == "href" {
					urlMap[attr.Val] = true
				}
			}
		}
	}
}

func cleanLinks(urlMap map[string]bool) (cleanedUrlMap map[string]bool) {
	for url := range urlMap {
		// Remove same-page element id links
		if strings.HasPrefix(url, "#") {
			delete(urlMap, url)
			continue
		}

		// Remove "//" prefix
		cleanUrl := strings.TrimPrefix(url, "//")

		// Pre-pend website domain name to keys with "/"
		if strings.HasPrefix(cleanUrl, "/") {
			cleanUrl = fmt.Sprintf("%s%s", START_URL, cleanUrl)
		}

		// Pre-pend HTTPS protocol if not exists
		if !strings.HasPrefix(cleanUrl, "http://") && !strings.HasPrefix(cleanUrl, "https://") {
			cleanUrl = fmt.Sprintf("%s%s", "https://", cleanUrl)
		}

		// Swap old URL for clean URL
		delete(urlMap, url)
		urlMap[cleanUrl] = true
	}

	return urlMap
}

func prettyPrint(v interface{}) (err error) {
	res, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		return
	}

	fmt.Println(string(res))

	return
}
