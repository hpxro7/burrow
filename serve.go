package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hpxro7/burrow/crawl"
)

type ServeOn chan chan string

const (
	updatePoolSize  = 120
	requestPoolSize = 20
	urlSinkSize     = 150
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

	crawledUrlSink, requestPool := CrawlMonitor(updatePoolSize, requestPoolSize, urlSinkSize)

	crawl.Through(crawl.UrlsUsingAnchor).BeginWith(seedUrls, crawledUrlSink)
	http.Handle("/geturl", ServeOn(requestPool))
	err = http.ListenAndServe("0.0.0.0:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func CrawlMonitor(crawlBufSize, requestBufSize, maxUrlSinkSize int) (crawled chan string, requests chan chan string) {
	crawled, requests = make(chan string, crawlBufSize), make(chan chan string, requestBufSize)
	go func() {
		var urls []string
		visited := make(map[string]bool)

		for {
			if len(urls) < maxUrlSinkSize {
				select {
				case next := <-crawled:
					if !visited[next] {
						urls = append(urls, next)
						visited[next] = true
						log.Println("Read into sink. New sink size: ", len(urls))
					}
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
