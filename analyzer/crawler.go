package analyzer

import (
	"log/slog"
	"slices"
)

func crawl(info *AnalysisData, status *parsingState) {
	linkStats := LinkStats{}
	if len(status.allLinks) == 0 {
		slog.Info("Nothing to crawl! No links found in the ", "site", info.ID())
		return
	}

	for k := range status.allLinks {
		if isAbsoluteUrl(k) {
			linkStats.ExternalLinkCount++
		} else {
			linkStats.InternalLinkCount++
		}
	}
	info.LinkStats = linkStats

	crawlForValidity(info, status)
}

func crawlForValidity(info *AnalysisData, status *parsingState) {
	invalidlinkchannel := make(chan LinkStatus)
	count := 0

	slog.Info("Starting crawling for links...", "site", info.ID(), "pending#", len(status.allLinks))
	for k := range status.allLinks {
		if !isAnchorLink(k) {

			checkUrl := k
			count++

			go func() {
				crawlUrl(checkUrl, info, invalidlinkchannel)
			}()
		}
	}

	for event := range invalidlinkchannel {
		count--

		if !event.IsValid {
			slog.Info("Invalid link found!", "url", event.Url, "status", event.StatusCode)
			info.LinkStats.InvalidLinkCount++
			info.LinkStats.InvalidLinks = append(info.LinkStats.InvalidLinks, event.Url)
		}

		if count <= 0 {
			close(invalidlinkchannel)
		}
	}

	// sort links so that similar links will be placed close together.
	slices.Sort(info.LinkStats.InvalidLinks)

	if info.LinkStats.InvalidLinkCount == 0 {
		slog.Info("No invalid links found!", "site", info.ID())
	}

	slog.Info("Finished crawling all links in the ", "site", info.ID())
}

func crawlUrl(url string, info *AnalysisData, c chan LinkStatus) {
	checkUrl, err := getFinalUrl(url, info.SourceUrl)
	isValid := false
	statusCode := 999
	if err == nil {
		isValid, statusCode = findUrlValidity(checkUrl)
	}

	c <- LinkStatus{Url: checkUrl, IsValid: isValid, StatusCode: statusCode}
}
