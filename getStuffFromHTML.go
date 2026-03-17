package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	res := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return ""
	}

	if h1 := doc.Find("h1").Text(); strings.TrimSpace(h1) != "" {
		return strings.TrimSpace(h1)
	}
	return strings.TrimSpace(doc.Find("h2").Text())

}

func getFirstParagraphFromHTML(html string) string {
	res := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return ""
	}

	if p := doc.Find("main").Find("p").First().Text(); strings.TrimSpace(p) != "" {
		return strings.TrimSpace(p)
	}
	return strings.TrimSpace(doc.Find("p").First().Text())
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	res := strings.NewReader(htmlBody)

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return nil, err
	}

	urls := []string{}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		parsedURL, err := url.Parse(href)
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(parsedURL)
		urls = append(urls, resolvedURL.String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	res := strings.NewReader(htmlBody)

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return nil, err
	}

	urls := []string{}

	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")

		parsedURL, err := url.Parse(src)
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(parsedURL)
		urls = append(urls, resolvedURL.String())
	})

	return urls, nil
}

func extractPageData(html, pageURL string) PageData {
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		parsedURL = &url.URL{}
	}

	outgoingLinks, _ := getURLsFromHTML(html, parsedURL)
	imageURLs, _ := getImagesFromHTML(html, parsedURL)

	return PageData{
		URL:            pageURL,
		Heading:        getHeadingFromHTML(html),
		FirstParagraph: getFirstParagraphFromHTML(html),
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
