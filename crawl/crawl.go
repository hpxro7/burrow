// Package crawl implements a simple library for crawling through linked document entities.
package crawl

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// Through encapsulates a process for collecting outbound resource identifiers within a given document
// body.
type Through func(body string) []string

var (
	anchorMatcher = regexp.MustCompile("<a href ?= ?\"(https?://.+?)\"")
)

// BeginWith sends found outbound URIs through crawled using urls as a seed.
func (through Through) BeginWith(urls []string, crawled chan<- string) {
	for _, url := range urls {
		go recurse(url, crawled, through)
	}
}

func recurse(url string, crawled chan<- string, through Through) {
	uriContents, err := getContents(url)
	if err == nil {
		newUrls := through(uriContents)
		for _, newUrl := range newUrls {
			crawled <- newUrl
			go recurse(newUrl, crawled, through)
		}
	}
}

func getContents(uri string) (content string, err error) {
	resp, err := http.Get(uri)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	content = string(contentBytes)
	return
}

// UrlsUsingAnchor returns all anchor hrefs found in a given document body.
func UrlsUsingAnchor(body string) (next []string) {
	found := anchorMatcher.FindAllStringSubmatch(body, -1)

	for _, submatch := range found {
		if len(submatch) > 1 {
			next = append(next, submatch[1:]...)
		}
	}
	return
}
