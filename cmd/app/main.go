/*****************************************************************************************************************/

//	@author		Michael Roberts

/*****************************************************************************************************************/

package main

/*****************************************************************************************************************/

import (
	"flag"
	"fmt"
	"time"

	"github.com/michealroberts/koroutine-web-crawler/pkg/crawler"
	"github.com/xlab/treeprint"
)

/*****************************************************************************************************************/

func addNodes(branch treeprint.Tree, node *crawler.URLNode) {
	if node == nil {
		return
	}

	nodeBranch := branch.AddBranch(node.URL)

	for _, child := range node.Links {
		addNodes(nodeBranch, child)
	}
}

/*****************************************************************************************************************/

func main() {
	fmt.Println("Starting the crawler...")

	domain := flag.String("domain", "https://koroutine.tech", "The domain to crawl")

	depth := flag.Int("depth", 3, "The maximum depth to crawl")

	flag.Parse()

	if *domain == "" {
		fmt.Println("No domain provided")
		return
	}

	if *depth < 1 {
		fmt.Println("Invalid depth")
		return
	}

	fmt.Println("Crawling domain:", *domain)

	fmt.Println("Crawling depth:", *depth)

	// Create a new crawler instance:
	crawler := crawler.New()

	// Start timing
	start := time.Now()

	rootNode, err := crawler.Crawl(*domain, *depth)

	if err != nil {
		fmt.Println(err)
		return
	}

	// End timing
	elapsed := time.Since(start)

	// Ideally, we would like this to be less than 100ms:
	fmt.Printf("Crawling took %v", elapsed)

	// Create a new treeprint "tree":
	tree := treeprint.New()

	addNodes(tree, rootNode)

	fmt.Println(tree.String())
}

/*****************************************************************************************************************/
