package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Define a custom User-Agent to mimic a standard web browser
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

// Create a reusable HTTP client with default headers
var client = &http.Client{}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
			break
		}
	}
	return
}

func crawl(targetURL string, ch chan string, chFinished chan bool) {
	fmt.Printf("Crawling: %s\n", targetURL)

	// Create a new GET request manually to add custom headers
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Printf("ERR: failed to create request for %s: %v\n", targetURL, err)
		chFinished <- true
		return
	}
	
	// Set the User-Agent header to bypass bot detection
	req.Header.Set("User-Agent", userAgent)

	// Perform the request using our custom client
	resp, err := client.Do(req)

	defer func() {
		chFinished <- true
	}()

	if err != nil {
		fmt.Printf("ERR: failed to crawl %s: %v\n", targetURL, err)
		return
	}

	b := resp.Body
	defer b.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERR: received non-OK status for %s: %d\n", targetURL, resp.StatusCode)
		return // Exit the crawl function on error status
	}

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()
		
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return
			}
			fmt.Printf("ERR: Tokenizing error on %s: %v\n", targetURL, z.Err())
			return

		case html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				ok, href := getHref(t)

				if ok {
					// Handle relative URLs properly (from previous fix)
					base, err := url.Parse(targetURL)
					if err != nil {
						continue
					}
					relative, err := url.Parse(href)
					if err != nil {
						continue
					}
					resolvedURL := base.ResolveReference(relative).String()

					if strings.HasPrefix(resolvedURL, "http://") || strings.HasPrefix(resolvedURL, "https://") {
						ch <- resolvedURL
					}
				}
			}
		}
	}
}

func main() {
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

    if len(seedUrls) == 0 {
        fmt.Println("No seed URLs provided. Using default: https://example.com")
        seedUrls = []string{"https://example.com"}
    }

	chUrls := make(chan string)
	chFinished := make(chan bool)

	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}

	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}

	fmt.Println("\nFound", len(foundUrls), "unique urls across seeded pages:")
	for url := range foundUrls {
		fmt.Println("-" + url)
	}
	close(chUrls)
}
