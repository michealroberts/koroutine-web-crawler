# koroutine-web-crawler

🕷️ A simple web-crawler in Go lang

This project is a simple web-crawler in Go lang. It is designed to be simple, fast, and efficient.

---

## Considerations

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

## Implementation

The crawler is implemented using a simple breadth-first search algorithm. It starts with a seed URL, fetches the HTML content, extracts all the URLs, and then fetches the HTML content of each URL. The process continues until the maximum depth is reached handling the myriad of edge cases along the way.

## Dependencies

I am not an expert in printing trees, so I have use the xlab/treeprint package to print the tree structure of the crawled URLs. The package is available at https://github.com/xlab/treeprint.