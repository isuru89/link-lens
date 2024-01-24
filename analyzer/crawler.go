package analyzer

import (
	"log/slog"
	"slices"
)

// Crawl crawls all given links in the given base url and returns statistics about
// the nature of links encountered. Such as, whether a link is internal, external or invalid.
// Also, it reports all invalid links found separately.
func (c *OneDepthCrawler) Crawl(baseUrl string, links map[string]bool) *LinkStats {
	linkStats := &LinkStats{}
	if len(links) == 0 {
		return linkStats
	}

	for link := range links {
		if isAbsoluteUrl(link) {
			linkStats.ExternalLinkCount++
		} else {
			linkStats.InternalLinkCount++
		}
	}

	crawlForValidity(baseUrl, linkStats, links)
	return linkStats
}

func crawlForValidity(baseUrl string, stats *LinkStats, links map[string]bool) {
	invalidLinkChannel := make(chan LinkStatus)
	count := 0

	slog.Info("Starting crawling for links...", "site", baseUrl, "pending#", len(links))
	for k := range links {
		if !isAnchorLink(k) {

			checkUrl := k
			count++

			go func() {
				crawlUrl(checkUrl, baseUrl, invalidLinkChannel)
			}()
		}
	}

	for event := range invalidLinkChannel {
		count--

		if !event.IsValid {
			slog.Info("Invalid link found!", "url", event.Url, "status", event.StatusCode)
			stats.InvalidLinkCount++
			stats.InvalidLinks = append(stats.InvalidLinks, event.Url)
		}

		if count <= 0 {
			close(invalidLinkChannel)
		}
	}

	// sort links so that similar links will be placed close together.
	slices.Sort(stats.InvalidLinks)

	if stats.InvalidLinkCount == 0 {
		slog.Info("No invalid links found!")
	}

	slog.Info("Finished crawling all links in the ", "site", baseUrl)
}

func crawlUrl(url, baseUrl string, c chan LinkStatus) {
	checkUrl, err := getFinalUrl(url, baseUrl)
	isValid := false
	statusCode := 999
	if err == nil {
		isValid, statusCode = findUrlValidity(checkUrl)
	}

	c <- LinkStatus{Url: checkUrl, IsValid: isValid, StatusCode: statusCode}
}
