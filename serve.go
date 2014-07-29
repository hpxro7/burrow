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

	body = Through(urlsUsingPage).BeginWith(seedUrls, nil)
	http.HandleFunc("/", serveUrl)
	http.ListenAndServe("0.0.0.0:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (through Through) BeginWith(urls []string, outbuf <-chan []string) string {
	return through(urls[0])[0]
}

func urlsUsingPage(body string) []string {
	return []string{"www.google.com"}
}

func readSeedUrls(filename string) (seedUrls []string, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err == nil {
		seedUrls = strings.Split(string(contents), "\n")
	}
	return
}

func serveUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s", body)
}
