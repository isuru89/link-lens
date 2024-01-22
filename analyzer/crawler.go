package analyzer

import (
	"log"
	"slices"
)

func crawl(info *AnalysisData, status *parsingState) {
	linkStats := LinkStats{}
	if len(status.allLinks) == 0 {
		log.Printf("Nothing to crawl! No links found in the site: %s.", info.SourceUrl)
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

	log.Printf("Starting crawling for #%d links...", len(status.allLinks))
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
			log.Printf("Invalid link found! url=%s, status=%d", event.Url, event.StatusCode)
			info.LinkStats.InvalidLinkCount++
			info.LinkStats.InvalidLinks = append(info.LinkStats.InvalidLinks, event.Url)
		}

		if count <= 0 {
			close(invalidlinkchannel)
		}
	}

	slices.Sort(info.LinkStats.InvalidLinks)
	log.Printf("Finished crawling all links in the site %s", info.ID())
}

func crawlUrl(url string, info *AnalysisData, c chan LinkStatus) {
	checkUrl := getFinalUrl(url, info.SourceUrl)
	isValid, statusCode := findUrlValidity(checkUrl)

	c <- LinkStatus{Url: checkUrl, IsValid: isValid, StatusCode: statusCode}
}
