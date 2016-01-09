package main

import (
	"fmt"
	"sync"
)

type Set struct {
	v   map[string]bool
	mux sync.Mutex
}

func (s Set) Add(url string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.v[url] = true
}

func (s Set) Get(url string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.v[url]
}

var fetched Set

type Fetcher interface {
	// Fetch 返回 URL 的 body 内容，并且将在这个页面上找到的 URL 放到一个 slice 中。
	Fetch(url string) (body string, urls []string, err error)
}

func CrawlWorker(url string, depth int, fetcher Fetcher, quit chan int) {
	if depth <= 0 || fetched.Get(url) {
		quit <- 0
		return
	}
	body, urls, err := fetcher.Fetch(url)
	fetched.Add(url)
	if err != nil {
		fmt.Println(err)
		quit <- 0
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	wait := make(chan int, len(urls))
	for _, nextUrl := range urls {
		go CrawlWorker(nextUrl, depth-1, fetcher, wait)
	}
	for _ = range urls {
		<-wait
	}
	quit <- 0
}

func Crawl(url string, depth int, fetcher Fetcher) {
	wait := make(chan int)
	go CrawlWorker(url, depth, fetcher, wait)
	<-wait
}

func main() {
	fetched = Set{v: make(map[string]bool)}
	Crawl("http://golang.org/", 4, fetcher)
}

type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
