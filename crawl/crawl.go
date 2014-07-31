package crawl

type Through func(body string) []string

func (through Through) BeginWith(urls []string, crawled chan<- []string) {
	for _, url := range urls {
		go func() {
			//TODO(hpxro7): Read http content-body from url
			for {
				newUrls := through(url)
				crawled <- newUrls
			}
		}()
	}
}

func UrlsUsingPage(body string) []string {
	return []string{"www.google.com", "www.yahoo.com"}
}
