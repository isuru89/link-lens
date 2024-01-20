package analyzer

import (
	"log"
	"sync"
)

func crawl(info *AnalysisData, status *parsingState) {
	linkStats := LinkStats{}
	if len(status.allLinks) == 0 {
		log.Printf("Nothing to crawl! No links found in the site: %s.", info.SourceUrl)
		return
	}

	for k := range status.allLinks {
		if isAbsoluteUrl(k) {
			linkStats.ExternalLinks++
		} else {
			linkStats.InternalLinks++
		}
	}
	info.LinkStats = linkStats

	crawlForValidity(info, status)
}

func crawlForValidity(info *AnalysisData, status *parsingState) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	log.Printf("Starting crawling for #%d links...", len(status.allLinks))
	for k := range status.allLinks {
		if !isAnchorLink(k) {
			wg.Add(1)

			checkUrl := k
			go func() {
				defer wg.Done()
				crawlUrl(checkUrl, info, &mutex)
			}()
		}
	}

	wg.Wait()

	log.Printf("Finished crawling all links in the site %s", info.ID())
}

func crawlUrl(url string, info *AnalysisData, mtx *sync.Mutex) {
	checkUrl := getFinalUrl(url, info.SourceUrl)
	isValid, statusCode := findUrlValidity(checkUrl)

	if !isValid {
		log.Printf("Inaccessible link found! %s [Status: %d]", checkUrl, statusCode)
		mtx.Lock()
		info.LinkStats.InvalidLinks++
		mtx.Unlock()
	} else {
		log.Printf("link ok! %s [Status: %d]", checkUrl, statusCode)
	}
}
