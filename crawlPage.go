package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
	counts             map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	// cfg.mu.Lock()
	// if len(cfg.pages) >= cfg.maxPages {
	// 	cfg.mu.Unlock()
	// 	cfg.wg.Done()
	// 	return
	// }
	// cfg.mu.Unlock()
	cfg.concurrencyControl <- struct{}{}

	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("%s is not a valid URL", rawCurrentURL)
		return
	}
	if cfg.baseURL.Host != current.Host {
		return
	}
	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("COULDNT NORMALIZE: %v", err)
		return
	}

	if isFirst := cfg.addPageVisit(normalizedCurrentURL); !isFirst {
		return
	}

	currentURLHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("COULDNT get HTML: %v", normalizedCurrentURL)
		return
	}
	fmt.Printf("\nGOT HTML from: %v", normalizedCurrentURL)

	cfg.mu.Lock()
	cfg.pages[normalizedCurrentURL] = extractPageData(currentURLHTML, rawCurrentURL)
	cfg.mu.Unlock()
	urls, err := getURLsFromHTML(currentURLHTML, cfg.baseURL)
	if err != nil {
		fmt.Printf("COULDNT get URLS FROM: %v", normalizedCurrentURL)
		return
	}
	for _, u := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(u)

	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// Check maxPages here, while the lock is held
	if len(cfg.pages) >= cfg.maxPages {
		return false
	}

	cfg.counts[normalizedURL]++
	if cfg.counts[normalizedURL] > 1 {
		return false
	}

	return true
}
