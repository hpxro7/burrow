package crawl

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Through func(body string) []string

var (
	anchorMatcher = regexp.MustCompile("<a href ?= ?\"(https?://.+?)\"")
)

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

func UrlsUsingAnchor(body string) (next []string) {
	found := anchorMatcher.FindAllStringSubmatch(body, -1)

	for _, submatch := range found {
		if len(submatch) > 1 {
			next = append(next, submatch[1:]...)
		}
	}
	return
}
