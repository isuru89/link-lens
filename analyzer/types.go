package analyzer

const (
	LoginForm = "LoginForm"
	Unknown   = "Unknown"
)

type LinkStatus struct {
	Url        string
	IsValid    bool
	StatusCode int
}

type LinkStats struct {
	InternalLinkCount int
	ExternalLinkCount int
	InvalidLinkCount  int
	InvalidLinks      []string
}

// Base interface for all possible crawling strategies.
type Crawler interface {
	Crawl(baseUrl string, links map[string]bool) *LinkStats
}

// Crawl only to a single level depth.
type OneDepthCrawler struct {
}

type AnalysisData struct {
	SourceUrl     string
	HtmlVersion   string
	Title         string
	HeadingsCount map[string]int
	LinkStats     LinkStats
	PageType      string
}

// Stores internal analysis and parsing status.
type parsingState struct {
	allLinks        map[string]bool
	currTag         string
	inputTypeCounts map[string]int
}

func NewAnalysis(url string) *AnalysisData {
	return &AnalysisData{
		SourceUrl:     url,
		HtmlVersion:   "5",
		HeadingsCount: map[string]int{},
		LinkStats:     LinkStats{},
	}
}
