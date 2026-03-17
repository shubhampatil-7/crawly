package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	args := os.Args[1:]

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	if len(args) < 3 {
		fmt.Println("not enough arguments provided")
		fmt.Println("usage: crawler <url> <maxConcurrency> <maxPages>")
		os.Exit(1)
	}
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("invalid Max Concurrency argument: %s", err)
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("invalid Max Pages argument: %s", err)
		os.Exit(1)
	}
	baseUrl, err := url.Parse(args[0])
	if err != nil {
		fmt.Printf("Couldn't Parse URL: %s", err)
		os.Exit(1)
	}

	cfg := &config{
		pages:              make(map[string]PageData),
		counts:             make(map[string]int),
		baseURL:            baseUrl,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	fmt.Printf("Starting crawl of: %s\n", baseUrl)

	cfg.wg.Add(1)
	go cfg.crawlPage(baseUrl.String())
	cfg.wg.Wait()
	fmt.Println("\nCrawling Done")
	err = writeJSONReport(cfg.pages, "report.json")
	if err != nil {
		fmt.Printf("couldn't write report: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Report written to report.json")
	os.Exit(0)

	os.Exit(0)
}

func getHTML(rawURL string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", rawURL, nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("got HTTP error: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")

	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("got non-HTTP response: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil

}
