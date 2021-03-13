package main

import (
	"fmt"
	"strings"
)

// Assumtions:
// - Prepended URL in cleanup step is always the Wikipedia start URL
// - HTTPS used exclusively

func TrimDoubleSlashPrefix(link Link) Link {
	link.Url = strings.TrimPrefix(link.Url, "//")
	return link
}

func PrefixWithMissingDomain(link Link) Link {
	if strings.HasPrefix(link.Url, "/") {
		link.Url = fmt.Sprintf("%s%s", START_URL, link.Url)
	}
	return link
}

func PrefixWithMissingProtocol(link Link) Link {
	if !strings.HasPrefix(link.Url, "http://") && !strings.HasPrefix(link.Url, "https://") {
		link.Url = fmt.Sprintf("%s%s", "https://", link.Url)
	}
	return link
}

var Cleaners = []CleanFunction{
	TrimDoubleSlashPrefix,
	PrefixWithMissingDomain,
	PrefixWithMissingProtocol,
}
