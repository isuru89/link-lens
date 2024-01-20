package analyzer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcatUrl(t *testing.T) {
	testcases := []struct {
		base     string
		path     string
		expected string
	}{
		{base: "", path: "", expected: ""},
		{base: "www.a.com", path: "x/y/z", expected: "www.a.com/x/y/z"},
		{base: "www.a.com", path: "/x/y/z", expected: "www.a.com/x/y/z"},
		{base: "www.a.com/", path: "/x/y/z", expected: "www.a.com/x/y/z"},
		{base: "www.a.com/", path: "", expected: "www.a.com/"},
		{base: "www.a.com", path: "", expected: "www.a.com"},
	}

	for _, v := range testcases {
		result := concatUrl(v.base, v.path)
		if result != v.expected {
			t.Fatalf("concatUrl() => Expected: %s, but recieved: %s", v.expected, result)
		}
	}
}

func TestGetFinalUrl(t *testing.T) {
	testcases := []struct {
		href      string
		url       string
		expected  string
		doesPanic bool
	}{
		{href: "", url: "", expected: "", doesPanic: true},
		{href: "", url: "https://www.a.com", expected: "https://www.a.com", doesPanic: false},
		{href: "", url: "http://www.a.com", expected: "http://www.a.com", doesPanic: false},
		{href: "", url: "ftp://www.a.com", expected: "", doesPanic: true},
		{href: "", url: "ws://www.a.com", expected: "", doesPanic: true},
		{href: "ftp://www.a.com", url: "https://www.a.com", expected: "", doesPanic: true},
		{href: "ws://www.a.com", url: "https://www.a.com", expected: "", doesPanic: true},
		{href: "http://www.a.com", url: "https://www.a.com", expected: "http://www.a.com", doesPanic: false},
		{href: "https://www.a.com", url: "https://www.b.com/x", expected: "https://www.a.com", doesPanic: false},
		{href: "https://www.b.com/y", url: "https://www.b.com/x", expected: "https://www.b.com/y", doesPanic: false},
		{href: "#anchor", url: "https://www.a.com", expected: "https://www.a.com#anchor", doesPanic: false},
		{href: "#anchor", url: "https://www.a.com/", expected: "https://www.a.com/#anchor", doesPanic: false},
		{href: "#anchor", url: "https://www.a.com/x", expected: "https://www.a.com/x#anchor", doesPanic: false},
		{href: "relurl1", url: "https://www.a.com/x", expected: "https://www.a.com/relurl1", doesPanic: false},
		{href: "relurl1", url: "https://www.a.com/x/y", expected: "https://www.a.com/x/relurl1", doesPanic: false},
		{href: "relurl1", url: "https://www.a.com", expected: "https://www.a.com/relurl1", doesPanic: false},
		{href: "/relurl1", url: "https://www.a.com/x/y", expected: "https://www.a.com/relurl1", doesPanic: false},
		{href: "/relurl1", url: "https://www.a.com", expected: "https://www.a.com/relurl1", doesPanic: false},
	}

	for _, v := range testcases {
		if v.doesPanic {
			assert.Panics(t, func() { getFinalUrl(v.href, v.url) },
				fmt.Sprintf("Code should panic! concatUrl('%s', '%s')", v.href, v.url))
		} else {
			result := getFinalUrl(v.href, v.url)
			if result != v.expected {
				t.Fatalf("concatUrl('%s', '%s') => Expected: %s, but recieved: %s", v.href, v.url, v.expected, result)
			}
		}
	}
}
