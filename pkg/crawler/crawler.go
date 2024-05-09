/*****************************************************************************************************************/

//	@author		Michael Roberts

/*****************************************************************************************************************/

package crawler

/*****************************************************************************************************************/

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	parse "github.com/michealroberts/koroutine-web-crawler/pkg/parsers"
)

/*****************************************************************************************************************/

type URLNode struct {
	URL   string     `json:"url"`
	Links []*URLNode `json:"links"`
}

/*****************************************************************************************************************/

type Crawler struct {
	baseDomain string
	visited    map[string]bool
	mu         sync.Mutex
	wg         sync.WaitGroup
	client     *http.Client
}

/*****************************************************************************************************************/

// NewCrawler creates a new instance of Crawler with an initialized HTTP client and visited map.
func New() *Crawler {
	return &Crawler{
		visited: make(map[string]bool),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

/*****************************************************************************************************************/

// Crawl starts the crawling process from a given URL up to a maximum depth.
func (c *Crawler) Crawl(startURL string, maxDepth int) (*URLNode, error) {
	parsedURL, err := url.Parse(startURL)

	if err != nil {
		return nil, err
	}

	c.baseDomain = parsedURL.Host

	root := &URLNode{URL: startURL}

	c.wg.Add(1)
	go c.crawlRecursive(startURL, root, 0, maxDepth)
	c.wg.Wait()

	return root, nil
}

/*****************************************************************************************************************/

func (c *Crawler) crawlRecursive(currentURL string, node *URLNode, depth int, maxDepth int) {
	defer c.wg.Done()

	if depth > maxDepth || c.hasVisited(currentURL) {
		return
	}

	c.markAsVisited(currentURL)

	links, err := c.fetchAndParse(currentURL)
	if err != nil {
		return
	}

	for _, link := range links {
		parsedLink, err := url.Parse(link)

		if err != nil || parsedLink.Host != c.baseDomain {
			continue
		}

		childNode := &URLNode{URL: link}

		node.Links = append(node.Links, childNode)

		c.wg.Add(1)

		go func() {
			c.crawlRecursive(link, childNode, depth+1, maxDepth)
		}()
	}
}

/*****************************************************************************************************************/

// isVisited checks if a URL has already been visited.
func (c *Crawler) hasVisited(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.visited[url]
}

/*****************************************************************************************************************/

// markAsVisited marks a URL as visited.
func (c *Crawler) markAsVisited(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.visited[url] = true
}

/*****************************************************************************************************************/

// fetchAndParse retrieves the HTML content from the specified URL and extracts links.
func (c *Crawler) fetchAndParse(urlStr string) ([]string, error) {
	resp, err := c.client.Get(urlStr)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK || resp.Header.Get("Content-Type") != "text/html" {
		return nil, fmt.Errorf("non-200 status code received: %d", resp.StatusCode)
	}

	return parse.AhrefsFromHTML(resp.Body, urlStr), nil
}

/*****************************************************************************************************************/
