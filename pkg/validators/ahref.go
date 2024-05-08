/*****************************************************************************************************************/

//	@author		Michael Roberts

/*****************************************************************************************************************/

package validate

/*****************************************************************************************************************/

import (
	"net/url"
)

/*****************************************************************************************************************/

// ValidateAHref validates an href URL against a base URL and returns the resolved URL.
func Ahref(base string, href string) (string, error) {
	// Parse the href URL:
	uri, err := url.Parse(href)

	// If there's an error parsing the URL, return it:
	if err != nil {
		return "", err
	}

	// Parse the base URL:
	baseURL, err := url.Parse(base)

	// If there's an error parsing the base URL, return it:
	if err != nil {
		return "", err
	}

	// If the href URL is absolute, return it:
	return baseURL.ResolveReference(uri).String(), nil
}

/*****************************************************************************************************************/
