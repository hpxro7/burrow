Burrow
======
An experiment in writing a document crawler in Go.

The API aims to be expressive yet succinct. Take for example the task of crawling through html documents using anchor hrefs:

```go
crawl.Through(urlsUsingAnchor).BeginWith(seedUrls, crawledUrlSink)
```

Note that the current example sever implementation does not persist crawled entities to disk but rather keeps a pool of urls
and polls and removes them as the reqeust multiplexer deems fit. Therefore if scalibility is a concern and you expect more than urlSinkSize concurrent requests I would recommend using an actual crawling engine.
