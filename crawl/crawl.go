package crawl

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Through func(body string) []string

func (through Through) BeginWith(urls []string, crawled chan<- string) {
	for _, url := range urls {
		go func() {
			uriContents, err := getContents(url)
			if err == nil {
				newUrls := through(uriContents)
				for _, newUrl := range newUrls {
					crawled <- newUrl
				}
			}
		}()
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

func UrlsUsingPage(body string) (urls []string) {
	return []string{body}
}
