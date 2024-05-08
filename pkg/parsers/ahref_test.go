/*****************************************************************************************************************/

//	@author		Michael Roberts

/*****************************************************************************************************************/

package parse

import (
	"io"
	"os"
	"strings"
	"testing"
)

/*****************************************************************************************************************/

func TestAhrefsFromHTML(t *testing.T) {
	// Load the HTML content from file
	fileContent, err := os.ReadFile("ahref_test.html")
	if err != nil {
		t.Fatalf("Failed to read HTML file: %v", err)
	}
	reader := io.NopCloser(strings.NewReader(string(fileContent)))

	// Call the function
	result := AhrefsFromHTML(reader, "http://base.com")

	// Expected results, considering how `validate.Ahref` might resolve URLs
	expected := []string{
		"http://example.com",
		"http://space.com",
		"http://base.com/relative/path",
		"http://valid.com/query?name=John%20Doe&status=active",
		"http://base.com",
		"http://valid.com/path",
	}

	// Check results
	if len(result) != len(expected) {
		t.Errorf("Expected %d hrefs, got %d", len(expected), len(result))
	}

	for i, href := range expected {
		if result[i] != href {
			t.Errorf("Expected href %s, got %s", href, result[i])
		}
	}
}

/*****************************************************************************************************************/
