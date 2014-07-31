package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hpxro7/picserve/crawl"
)

type ServeOn chan chan string

const (
	updatePoolSize  = 10
	requestPoolSize = 5
	urlPoolSize     = 20
)

var (
	seedFilename = flag.String("seed_file", "", "The file containing a url on each line to be used as crawling seeds.")
)

func init() {
	flag.Parse()
}

func main() {
	if *seedFilename == "" {
		log.Fatal("Seed filename must be specified")
		os.Exit(2)
	}

	seedUrls, err := readSeedUrls(*seedFilename)
	if err != nil {
		log.Fatal("Could not read file")
		os.Exit(1)
	}

	crawlPool, requestPool := CrawlMonitor(updatePoolSize, requestPoolSize, urlPoolSize)

	crawl.Through(crawl.UrlsUsingPage).BeginWith(seedUrls, crawlPool)
	http.Handle("/geturl", ServeOn(requestPool))
	err = http.ListenAndServe("0.0.0.0:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func CrawlMonitor(crawlBufSize, requestBufSize, maxUrlPoolSize int) (crawls chan []string, requests chan chan string) {
	crawls, requests = make(chan []string, crawlBufSize), make(chan chan string, requestBufSize)
	go func() {
		var urls []string
		for {
			if len(urls) < maxUrlPoolSize {
				next := <-crawls
				urls = append(urls, next...)
				log.Println("Saved urls:", next)
			}

			if len(urls) >= 0 {
				select {
				case req := <-requests:
					req <- urls[0]
					urls = urls[1:]
				default:
				}
			}
		}
	}()
	return
}

func (requestPool ServeOn) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	urlRequest := make(chan string)
	requestPool <- urlRequest
	fmt.Fprintf(w, "You got back a: %s", <-urlRequest)
}

func readSeedUrls(filename string) (seedUrls []string, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err == nil {
		seedUrls = strings.Split(string(contents), "\n")
	}
	return
}
