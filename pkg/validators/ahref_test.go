/*****************************************************************************************************************/

//	@author		Michael Roberts

/*****************************************************************************************************************/

package validate

/*****************************************************************************************************************/

import "testing"

/*****************************************************************************************************************/

func TestValidateAHref(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		href     string
		expected string
	}{
		{
			name:     "Basic relative path",
			base:     "http://example.com/dir/",
			href:     "subdir/page.html",
			expected: "http://example.com/dir/subdir/page.html",
		},
		{
			name:     "Absolute href overrides base",
			base:     "http://example.com/dir/",
			href:     "https://another.com/page",
			expected: "https://another.com/page",
		},
		{
			name:     "Href with query",
			base:     "http://example.com/",
			href:     "page?query=123",
			expected: "http://example.com/page?query=123",
		},
		{
			name:     "Href with anchor",
			base:     "http://example.com/",
			href:     "page#anchor",
			expected: "http://example.com/page#anchor",
		},
		{
			name:     "Going up one directory",
			base:     "http://example.com/dir/",
			href:     "../updir/page.html",
			expected: "http://example.com/updir/page.html",
		},
		{
			name:     "Empty href returns base",
			base:     "http://example.com/dir/",
			href:     "",
			expected: "http://example.com/dir/",
		},
		{
			name:     "Special characters",
			base:     "http://example.com/",
			href:     "dir/page?.html",
			expected: "http://example.com/dir/page?.html",
		},
		{
			name:     "Trailing slash and file",
			base:     "http://example.com/dir/",
			href:     "file.txt",
			expected: "http://example.com/dir/file.txt",
		},
		{
			name:     "Whitespam Trim",
			base:     "http://example.com/",
			href:     "   http://example.com/   ",
			expected: "http://example.com/",
		},
		{
			name:     "Javascript void",
			base:     "http://example.com/",
			href:     "javascript:void(0);",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uri, err := Ahref(tc.base, tc.href)

			if err != nil {
				t.Errorf("Test %s failed: %s", tc.name, err)
			}

			if uri != tc.expected {
				t.Errorf("Test %s failed: expected %s, got %s", tc.name, tc.expected, uri)
			}
		})
	}
}

/*****************************************************************************************************************/
