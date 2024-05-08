/*****************************************************************************************************************/

//	@author		Michael Roberts

/*****************************************************************************************************************/

package parse

/*****************************************************************************************************************/

import (
	"io"
	"strings"

	validate "github.com/michealroberts/koroutine-web-crawler/pkg/validators"
	"golang.org/x/net/html"
)

/*****************************************************************************************************************/

// Extracts all anchor tags from an HTML document and returns the href URLs.
func AhrefsFromHTML(body io.ReadCloser, base string) []string {
	// Placeholder for extracted ahrefs:
	var ahrefs []string

	// Create an HTML tokenizer to parse the content:
	tokenizer := html.NewTokenizer(body)

	// Parse the HTML content and extract ahrefs from anchor tags:
	for {
		// Get the next token type:
		tokenType := tokenizer.Next()

		// If it's an error token, we've reached the end of the document:
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()

		if tokenType == html.StartTagToken && token.Data == "a" {
			for _, a := range token.Attr {
				if a.Key == "href" {
					resolvedAhref, err := validate.Ahref(base, a.Val)

					if err != nil {
						continue
					}

					if strings.HasPrefix(resolvedAhref, "http") {
						ahrefs = append(ahrefs, resolvedAhref)
					}
				}
			}
		}
	}

	return ahrefs
}

/*****************************************************************************************************************/
