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
	Root       *URLNode
	baseDomain string
	visited    map[string]bool
	mu         sync.Mutex
	wg         sync.WaitGroup
	client     *http.Client
	stream     chan *URLNode // channel for streaming URL nodes
}

/*****************************************************************************************************************/

// NewCrawler creates a new instance of Crawler with an initialized HTTP client and visited map.
func New() *Crawler {
	root := &URLNode{} // Initialize with a root node if necessary

	return &Crawler{
		Root:    root,
		visited: make(map[string]bool),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		stream: make(chan *URLNode, 100), // buffered channel to avoid blocking
	}
}

/*****************************************************************************************************************/

// Stream provides the channel to receive streamed results
func (c *Crawler) Stream() <-chan *URLNode {
	return c.stream
}

/*****************************************************************************************************************/

// Crawl starts the crawling process from a given URL up to a maximum depth.
func (c *Crawler) Crawl(startURL string, maxDepth int) (*URLNode, error) {
	defer close(c.stream) // Ensure the channel is closed when done

	parsedURL, err := url.Parse(startURL)

	if err != nil {
		return nil, err
	}

	c.baseDomain = parsedURL.Host

	root := &URLNode{URL: startURL}

	c.Root = root

	c.wg.Add(1)
	go c.crawlRecursive(startURL, c.Root, 0, maxDepth)
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

		c.mu.Lock()
		node.Links = append(node.Links, childNode)
		c.mu.Unlock()

		c.stream <- childNode // send childNode to the channel

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
