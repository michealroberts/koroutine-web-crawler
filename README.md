# koroutine-web-crawler Part I

🕷️ A simple web-crawler in Go lang

This project is a simple web-crawler in Go lang. It is designed to be simple, fast, and efficient.

---

## Usage

```go
package main

// Create a new crawler instance:
crawler := crawler.New()

// Stream the crawled output as we receive it:
go func() {
  for node := range crawler.Stream() {
		fmt.Printf("Crawled URL: %s with %d links\n", node.URL, len(node.Links))
  }
}()

// Crawl some poor website to a maximum depth of 3:
rootNode, err := crawler.Crawl("https://example.com", 3)
```

The root structure contains a tree-like structure of recursively parsed ahrefs and their validated URL link, for example:

```go
type URLNode struct {
	URL   string
	Links []*URLNode
}

rootNode := &URLNode{
  URL: "https://example.com",
  Links: []*URLNode{
    {
      URL: "https://example.com/about",
      Links: []*URLNode{
        {
          URL:   "https://example.com/about/team",
          Links: nil,
        },
        {
          URL:   "https://example.com/about/careers",
          Links: nil,
        },
        {
          URL:   "https://example.com/contact",
          Links: nil,
        },
      },
    },
    {
      URL:   "https://example.com/contact",
      Links: nil,
    },
  },
}
```

The crawler will only crawl to the maximum recursion depth provided, avoiding duplicates within an individual node, but may contain overlapping URLs in different nodes.

## API

## Example Commands

To run a development container, the user can use the following command:

```bash
make dev
```

Or open the workspace in Visual Studio Code and use the development container. See the `.devcontainer` folder for more information.

The crawler is designed to be simple and easy to use. 

The user can also provide the seed URL and the maximum depth as command-line arguments.

```bash
go run ./cmd/app/main.go -domain=https://example.com -depth=3
```

This will crawl the website `https://example.com` to a maximum depth of 3, and print the tree structure of the crawled URLs.

There is also the ability to run a server that listens for incoming requests and returns the tree structure of the crawled URLs, as a server sent stream of data, adhering to the SSE standard.

To run the Gin API, simply use: 

```bash
go run ./cmd/api/main.go
```

To then test the SSE route, the user can use the following command:

```bash
curl -N "http://localhost:8080/crawl?domain=https://example.com&depth=2"
```

This will stream the output, gathering links from the website `https://example.com` to a maximum depth of 2 and streaming back to the client.

## Local Development Setup

The user can setup the development environment by either using the Dockerfile provided, or by utilising the devcontainer in Visual Studio Code.

To use the Dockerfile in interactive terminal mode, the user can run the following commands:

```bash
make dev
```

Under the hood, the `make dev` command will build the Docker image and run the container in interactive mode, e.g.,

```bash
export CONTAINER_NAME="koroutine/web-crawler"

docker build --target development --tag $(CONTAINER_NAME):dev .

docker run -it --rm -v $(pwd):/app -w /app koroutine-web-crawler:latest
```

The user can then run the tests, build the project, and run the crawler, etc all from inside a development container:

## VSCode Development Container

For the latter, the user will be prompted to reload the workspace in the devcontainer. The user can then run the tests, build the project, and run the crawler, etc.

## Considerations & Thought Process

At first glance, crawling web-URLs is simple. However, we need to be cautious. Most websites will have sophisticated anti-crawling strategies in place, urls can be broken, and the crawler can get stuck in a loop. We want to crawl a website, recursively, but we do not want our dinner to get cold and be stuck waiting for hours. 

- Domain Restriction: 

The crawler extracts the domain from the start URL and ensures that subsequent requests only follow links within this domain, strictly adhering to the given domain without crawling subdomains or external domains, for example if the seed root URL is `https://example.com`, the crawler will only follow links that start with `https://example.com`, not `https://sub.example.com` or `https://another.com`.

- Maximum Depth (Recursion Pit):

The crawler has a maximum depth limit to prevent it from crawling infinitely. The maximum depth is set to 3 by default, but it can be changed by the user. It is advised that the user uses a sensible value to prevent the crawler from getting stuck in a loop, or crawling the entire internet.

- Handling !2** Status Codes:

The crawler needs to handle the case where the response code is not a 200 OK. In such cases, the crawler should not follow the link and should continue with the next link. Although any 2** status code is considered a success, the crawler should not follow the link if the status code is not 200 OK as this is the HTML specification for a successful response.

- Ahref Validation

We need to ensure that the ahrefs are validated to some standard to ensure that the crawler does not return broken or invalid links, or links that are not actually URLs.

## Testing Strategy

The testing strategy focused on these core unit tests:

- Testing the URL extraction from the HTML content
- Testing the URL normalization
- Testing the URL filtering based on the domain
- Testing the URL filtering based on the status code
- Testing the URL filtering based on the depth
- Testing for user input related edge cases

My philosophy on tests is to begin with what makes sense, and as you run into issues, add more tests to cover more edge cases cases and scenarios as we find them.

## Implementation

The crawler is implemented using a simple breadth-first search algorithm. It starts with a seed URL, fetches the HTML content, extracts all the URLs, and then fetches the HTML content of each URL. The process continues until the maximum depth is reached handling the myriad of edge cases along the way.

## Dependencies

This project is designed to be dependency light, and only uses the following dependencies:

- github.com/xlab/treeprint

I am not an expert in printing trees, so I have use the xlab/treeprint package to print the tree structure of the crawled URLs. The package is available at https://github.com/xlab/treeprint.

- github.com/stretchr/testify/assert

Because nobody likes to reinvent the testing wheel.