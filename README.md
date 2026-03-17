# Crawly

A concurrent web crawler written in Go. Given a starting URL, it recursively crawls all pages within the same domain and outputs a structured JSON report.

## Features

- Concurrent crawling with configurable parallelism
- Stays within the origin domain (no external links followed)
- Deduplicates visited URLs
- Respects a configurable page limit
- Extracts per-page data: heading, first paragraph, outgoing links, and images
- Outputs a sorted JSON report

## Building

```
go build -o crawly .
./crawly https://example.com 5 100
```
## Usage

```
go run . <url> <maxConcurrency> <maxPages>
```

| Argument         | Description                                      |
| ---------------- | ------------------------------------------------ |
| `url`            | The starting URL to crawl                        |
| `maxConcurrency` | Maximum number of pages fetched in parallel      |
| `maxPages`       | Maximum number of pages to crawl before stopping |

### Example

```
go run . https://example.com 5 100
```

This crawls `https://example.com` with up to 5 concurrent workers, stopping after 100 pages. Results are written to `report.json`.

## Output

`report.json` contains a sorted array of page objects:

```json
[
  {
    "url": "https://example.com/about",
    "heading": "About Us",
    "first_paragraph": "We are a company that...",
    "outgoing_links": ["https://example.com/contact", "..."],
    "image_urls": ["https://example.com/logo.png", "..."]
  }
]
```

## Project Structure

| File                   | Description                                            |
| ---------------------- | ------------------------------------------------------ |
| `main.go`              | Entry point, argument parsing, crawl orchestration     |
| `crawlPage.go`         | Core crawl logic, concurrency control, `config` struct |
| `getStuffFromHTML.go`  | HTML parsing: links, images, headings, paragraphs      |
| `extract_page_data.go` | Assembles `PageData` from a raw HTML string            |
| `normalize_url.go`     | Normalizes URLs (lowercase, strip trailing slash)      |
| `jsonreport.go`        | Writes the sorted JSON report to disk                  |

## Dependencies

- [goquery](https://github.com/PuerkitoBio/goquery) — HTML parsing and querying


