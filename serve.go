package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Through func(body string) []string
type Request chan chan string

var body string

func main() {
	seedFilename := flag.String("seed", "",
		"The file containing a url on each line to be used as crawling seeds.")
	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	seedUrls, err := readSeedUrls(*seedFilename)
	if err != nil {
		log.Fatal("Could not read file")
		os.Exit(1)
	}

	crawls, requests := CrawlMonitor()

	Through(urlsUsingPage).BeginWith(seedUrls, crawls)
	http.Handle("/geturl", Request(requests))
	http.ListenAndServe("0.0.0.0:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func CrawlMonitor() (crawls chan []string, requests chan chan string) {
	crawls, requests = make(chan []string, 20), make(chan chan string, 5)
	go func() {
		urls := make([]string, 0, 10)
		for {
			select {
			case nextList := <-crawls:
				log.Println("Added to urlList: ", nextList)
				urls = append(urls, nextList...)
			case req := <-requests:
				req <- urls[0]
				urls = urls[1:]
				log.Println("New urlList: ", urls)
			}
		}
	}()
	return
}

func (requestPool Request) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	urlReq := make(chan string)
	requestPool <- urlReq
	log.Println("Request:\n", req)
	fmt.Fprintf(w, "You got back a: %s", <-urlReq)
}

func (through Through) BeginWith(urls []string, crawled chan<- []string) {
	for _, url := range urls {
		go func() {
			//TODO(hpxro7): Read http content-body from url
			dataCrawled := through(url)
			fmt.Println(dataCrawled)
			crawled <- dataCrawled
		}()
	}
}

func urlsUsingPage(body string) []string {
	return []string{"www.google.com", "www.yahoo.com"}
}

func readSeedUrls(filename string) (seedUrls []string, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err == nil {
		seedUrls = strings.Split(string(contents), "\n")
	}
	return
}
