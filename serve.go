package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hpxro7/picserve/crawl"
)

type ServeOn chan chan string

const (
	updatePoolSize  = 120
	requestPoolSize = 20
	urlPoolSize     = 150
)

var (
	seedFilename = flag.String("seed_file", "", "The file containing a url on each line to be used as crawling seeds.")
)

func init() {
	flag.Parse()
}

func main() {
	if *seedFilename == "" {
		flag.PrintDefaults()
		log.Fatal("Seed filename must be specified")
	}

	seedUrls, err := readSeedUrls(*seedFilename)
	if err != nil {
		log.Fatal("Could not read file")
	}

	crawlPool, requestPool := CrawlMonitor(updatePoolSize, requestPoolSize, urlPoolSize)

	crawl.Through(crawl.UrlsUsingAnchor).BeginWith(seedUrls, crawlPool)
	http.Handle("/geturl", ServeOn(requestPool))
	err = http.ListenAndServe("0.0.0.0:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func CrawlMonitor(crawlBufSize, requestBufSize, maxUrlPoolSize int) (crawls chan string, requests chan chan string) {
	crawls, requests = make(chan string, crawlBufSize), make(chan chan string, requestBufSize)
	go func() {
		var urls []string
		for {
			if len(urls) < maxUrlPoolSize {
				select {
				case next := <-crawls:
					urls = append(urls, next)
					log.Println("Read into pool. New pool size: ", len(urls))
				default:
				}
			}

			if len(urls) > 0 {
				select {
				case req := <-requests:
					log.Println("Serving request...")
					req <- urls[0]
					urls = urls[1:]
					log.Println("...request served!")
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
	fmt.Fprint(w, <-urlRequest)
}

func readSeedUrls(filename string) (seedUrls []string, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err == nil {
		seedUrls = strings.Split(string(contents), "\n")
	}
	return
}
