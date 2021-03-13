package main

import (
	"fmt"
	"strings"
)

type Link struct {
	Url    string
	Domain string
}

type Links map[Link]bool
type CleanFunction func(Link) Link

func (l Link) Clean() Link {
	for _, currentCleanFunc := range Cleaners {
		l = currentCleanFunc(l)
	}

	return l
}

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
		linkSlice = append(linkSlice, fmt.Sprintf("{ domain: \"%s\", url: \"%s\" }", link.Domain, link.Url))
	}

	return strings.Join(linkSlice, "\n")
}

func (links Links) AddDomain(domain string) Links {
	newLinks := make(Links)

	for link := range links {
		link.Domain = domain
		newLinks[link] = true
	}

	return newLinks
}

func NewLink(url string) Link {
	return Link{
		Url: url,
	}
}
