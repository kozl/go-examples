package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type Cache struct {
	data map[string]string
	mux  sync.Mutex
}

func (c *Cache) Write(url, body string) {
	c.mux.Lock()
	c.data[url] = body
	c.mux.Unlock()
}

func (c *Cache) Check(url string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.data[url]; ok {
		return true
	}
	return false
}

var wg sync.WaitGroup

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, cache *Cache) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	defer wg.Done()
	if depth <= 0 {
		return
	}
	if !cache.Check(url) {
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		cache.Write(url, body)
		for _, u := range urls {
			wg.Add(1)
			go Crawl(u, depth-1, fetcher, cache)
		}
		return
	}
}

func main() {
	c := Cache{data: make(map[string]string)}
	wg.Add(1)
	Crawl("http://golang.org/", 4, fetcher, &c)
	wg.Wait()
	for url, body := range c.data {
		fmt.Printf("url: %v \nbody: %v \n --- \n", url, body)
	}
}

// fakeFetcher is Fetcher that returns canned results.
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

// fetcher is a populated fakeFetcher.
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
	"http://golang.org/cmd/": &fakeResult{
		"Go command",
		[]string{},
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
