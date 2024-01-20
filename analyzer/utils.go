package analyzer

import (
	"net/http"
	"regexp"
	"strings"
)

var baseUrlRegex = regexp.MustCompile(`(?i)(https?://[^/]+)/?`)

// concatUrl returns a valid concatenated url by combining both
// baseUrl and path parameters mixing forward slashes (/) correctly.
func concatUrl(baseUrl, path string) string {
	if path == "" {
		return baseUrl
	}

	startsWithHash := strings.HasPrefix(path, "/")
	endsWithHash := strings.HasSuffix(baseUrl, "/")

	if endsWithHash && startsWithHash {
		return baseUrl + strings.TrimPrefix(path, "/")
	} else if !endsWithHash && !startsWithHash {
		return baseUrl + "/" + path
	} else {
		return baseUrl + path
	}
}

// isAbsoluteUrl returns true if given href indicates a full absolute url.
func isAbsoluteUrl(href string) bool {
	return strings.Index(href, "://") >= 0
}

// isRelativeUrl returns true if the given href is a relative url
// with respective to the base url. If the url starts with a forward
// slash, then it indicates a relative path to the base url.
// There is another variation of relative url which does not handle by this method.
func isRelativeUrl(href string) bool {
	return strings.HasPrefix(href, "/")
}

// isAnchorLink returns true if the given href is a anchor link.
// Anchor links usually starts with a hash.
func isAnchorLink(href string) bool {
	return strings.HasPrefix(href, "#")
}

// getFinalUrl returns the final absolute url we need to fetch or check.
// This modifies the href as necessary with the source url analyzing.
func getFinalUrl(href, sourceUrl string) string {
	if sourceUrl == "" {
		panic("source URL cannot be empty!")
	}

	baseUrlParts := baseUrlRegex.FindStringSubmatch(sourceUrl)
	if baseUrlParts == nil {
		panic("Unable to find valid base url! Either its unsupported portocol or malformed url!")
	}

	if href == "" {
		return sourceUrl
	}

	if isAbsoluteUrl(href) {
		// check for valid protocol
		hrefParts := baseUrlRegex.FindStringSubmatch(href)
		if hrefParts == nil {
			panic("Unsupported href! Either its unsupported portocol or malformed url!")
		}
		return href
	} else if isAnchorLink(href) {
		return sourceUrl + href
	}

	baseUrl := baseUrlParts[1]

	if isRelativeUrl(href) {
		return concatUrl(baseUrl, href)
	}

	pos := strings.LastIndex(sourceUrl, "/")
	dpos := strings.LastIndex(sourceUrl, "//")
	if pos > 0 && dpos < pos-1 {
		return sourceUrl[:pos+1] + href
	} else if pos > 0 && dpos == pos-1 {
		return concatUrl(sourceUrl, href)
	}
	return href
}

// FindUrlValidity returns true if this given link is a valid one or not
// by checking whether it returns a 2xx response.
// Note: This method does not strictly check the content-type.
func FindUrlValidity(checkUrl string) (bool, int) {
	resp, err := http.Get(checkUrl)
	if err != nil {
		// we swallow the error, because caller cares only about status
		return false, 999
	}
	defer resp.Body.Close()

	// we still dont care sites returning html content with with status code >=400
	// e.g. Nginx 404/5xx
	if resp.StatusCode >= 300 {
		return false, resp.StatusCode
	}
	return true, resp.StatusCode
}
