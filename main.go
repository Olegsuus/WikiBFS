package main

import (
	"bufio"
	"container/list"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	GlobalTrace []string
	logger      *log.Logger
)

func initLogger() {
	file, err := os.OpenFile("crawler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	logger = log.New(file, "", log.LstdFlags)
}

func fetchLinks(pageURL string) ([]string, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Printf("Non-OK HTTP status: %d for URL %s", resp.StatusCode, pageURL)
		return nil, fmt.Errorf("non-OK HTTP status: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, nil
}

func UpgradeURL(links *[]string) []string {
	for i := range *links {
		(*links)[i], _ = readableURL(&(*links)[i])
		(*links)[i] = "https://ru.wikipedia.org" + (*links)[i]
	}
	return *links
}

func linkFilter(links []string) []string {
	filteredLinks := []string{}
	for _, link := range links {
		if strings.HasPrefix(link, "https://ru.wikipedia.org/wiki/") &&
			!strings.HasSuffix(link, ".svg") &&
			!strings.HasSuffix(link, ".jpg") &&
			!strings.HasPrefix(link, "https://ru.wikipedia.org/wiki/Википедия:") &&
			!strings.HasPrefix(link, "https://ru.wikipedia.org/wiki/Служебная:") &&
			!strings.HasPrefix(link, "https://ru.wikipedia.org/wiki/Категория:") &&
			link != "https://ru.wikipedia.org" {
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}

func bfs(startURL, targetURL string, traceLimit int) bool {
	queue := list.New()
	queue.PushBack([]string{startURL})
	visited := make(map[string]bool)
	visited[startURL] = true
	logger.Printf("Visited: %s", startURL)

	for queue.Len() > 0 {
		path := queue.Remove(queue.Front()).([]string)
		currentURL := path[len(path)-1]

		if len(path) > traceLimit {
			continue
		}

		if currentURL == targetURL {
			GlobalTrace = path
			return true
		}

		links, err := fetchLinks(currentURL)
		if err != nil {
			logger.Printf("Failed to fetch links from %s: %v", currentURL, err)
			continue
		}

		links = UpgradeURL(&links)
		links = linkFilter(links)

		for _, link := range links {
			if !visited[link] {
				visited[link] = true
				logger.Printf("Visited: %s", link)
				newPath := append([]string{}, path...)
				newPath = append(newPath, link)
				queue.PushBack(newPath)
			}
		}
	}

	return false
}

func main() {
	initLogger()

	startURL := readURLFromInput("Введите стартовую страницу: ")
	targetURL := readURLFromInput("Введите искомую страницу: ")
	traceLimit := 3

	if bfs(startURL, targetURL, traceLimit) {
		fmt.Println("Target found:", GlobalTrace)
		for i := 1; i < len(GlobalTrace); i++ {
			currentURL := GlobalTrace[i-1]
			targetURL = GlobalTrace[i]
			text, err := fetchParagraphWithLink(currentURL, targetURL)
			if err != nil {
				logger.Printf("Error fetching paragraph with link from %s to %s: %v", currentURL, targetURL, err)
				continue
			}
			fmt.Printf("%d------------------------\n%s\n%s\n", i, text, targetURL)
		}
	} else {
		fmt.Println("Target not found.")
	}
}

func readableURL(encodedURL *string) (string, error) {
	u, err := url.Parse(*encodedURL)
	if err != nil {
		return "", err
	}

	decodedPath, err := url.QueryUnescape(u.Path)
	if err != nil {
		return "", err
	}

	u.Path = decodedPath

	return u.Path, nil
}

func fetchParagraphWithLink(currentURL, targetURL string) (string, error) {
	resp, err := http.Get(currentURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	parsedTargetURL, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}

	var paragraphText string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			var pTextBuilder strings.Builder
			var containsLink bool
			var traverse func(*html.Node)
			traverse = func(c *html.Node) {
				if c.Type == html.ElementNode && c.Data == "a" {
					for _, attr := range c.Attr {
						if attr.Key == "href" {
							linkURL, err := url.Parse(attr.Val)
							if err == nil && linkURL.Path == parsedTargetURL.Path {
								containsLink = true
							}
						}
					}
				}
				if c.Type == html.TextNode {
					pTextBuilder.WriteString(c.Data)
				}
				for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
					traverse(gc)
				}
			}
			traverse(n)
			if containsLink {
				paragraphText = pTextBuilder.String()
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if paragraphText == "" {
		return "", fmt.Errorf("paragraph with link not found")
	}

	return paragraphText, nil
}

func readURLFromInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	url, _ := reader.ReadString('\n')
	return strings.TrimSpace(url)
}
