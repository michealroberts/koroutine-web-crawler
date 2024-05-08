/*****************************************************************************************************************/

//	@author		Michael Roberts

/*****************************************************************************************************************/

package crawler

/*****************************************************************************************************************/

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

/*****************************************************************************************************************/

func TestCrawlerInitialization(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
	assert.IsType(t, &http.Client{}, c.client)
	assert.NotNil(t, c.visited)
}

/*****************************************************************************************************************/

func TestCrawlValidStartURL(t *testing.T) {
	// Activate the http mock for the default HTTP client used globally
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Setup the base URL and the HTML content it should return
	baseURL := "https://koroutine.tech"

	httpmock.RegisterResponder("GET", baseURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `<html><head></head><body><a href="https://koroutine.tech/page1">Page 1</a></body></html>`)
			resp.Header.Add("Content-Type", "text/html")
			return resp, nil
		})

	// Create a new instance of the crawler
	c := New()
	// Perform the crawl operation starting from the base URL
	rootNode, err := c.Crawl(baseURL, 1)

	// Ensure no errors occurred
	assert.NoError(t, err)
	// Ensure a rootNode is returned
	assert.NotNil(t, rootNode)
	// The rootNode URL should be the baseURL
	assert.Equal(t, baseURL, rootNode.URL)
	// There should be exactly one link discovered
	assert.Len(t, rootNode.Links, 1)
	// The link discovered should be as expected
	assert.Equal(t, "https://koroutine.tech/page1", rootNode.Links[0].URL)
}

/*****************************************************************************************************************/

func TestCrawlIgnoreExternalLinks(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	baseURL := "https://koroutine.tech"

	httpmock.RegisterResponder("GET", baseURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `<a href="https://external.com">External</a>`)
			resp.Header.Add("Content-Type", "text/html")
			return resp, nil
		})

	c := New()
	root, err := c.Crawl(baseURL, 1)

	assert.NoError(t, err)
	assert.NotNil(t, root)
	assert.Empty(t, root.Links)
}

/*****************************************************************************************************************/

func TestCrawlNon200StatusCode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testURL := "https://koroutine.tech"
	httpmock.RegisterResponder("GET", testURL,
		httpmock.NewStringResponder(404, ""))

	c := New()
	root, err := c.Crawl(testURL, 1)

	assert.NoError(t, err)
	assert.NotNil(t, root)
	assert.Empty(t, root.Links)
}

/*****************************************************************************************************************/

func TestCrawlerConcurrency(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	baseURL := "https://koroutine.tech"
	link1 := "https://koroutine.tech/page1"
	link2 := "https://koroutine.tech/page2"

	httpmock.RegisterResponder("GET", baseURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `<a href="/page1">Page 1</a><a href="/page2">Page 2</a>`)
			resp.Header.Add("Content-Type", "text/html")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", link1,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Add("Content-Type", "text/html")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", link2,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Add("Content-Type", "text/html")
			return resp, nil
		})

	c := New()
	root, err := c.Crawl(baseURL, 2)

	assert.NoError(t, err)
	assert.NotNil(t, root)
	assert.Len(t, root.Links, 2)
	assert.Equal(t, link1, root.Links[0].URL)
	assert.Equal(t, link2, root.Links[1].URL)
}

/*****************************************************************************************************************/

func BenchmarkCrawler(b *testing.B) {
	c := New()

	for i := 0; i < b.N; i++ {
		c.Crawl("https:/koroutine.tech", 2)
	}
}

/*****************************************************************************************************************/
