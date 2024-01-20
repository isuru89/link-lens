package analyzer

const (
	LoginForm = "LoginForm"
	Unknown   = "Unknown"
)

type LinkStats struct {
	InternalLinks int
	ExternalLinks int
	InvalidLinks  int
}

type AnalysisData struct {
	SourceUrl     string
	HtmlVersion   string
	Title         string
	HeadingsCount map[string]int
	LinkStats     LinkStats
	PageType      string
}

type parsingState struct {
	allLinks        map[string]bool
	currTag         string
	inputTypeCounts map[string]int
}

func (a *AnalysisData) ID() string {
	return a.SourceUrl
}

func NewAnalysis(url string) *AnalysisData {
	return &AnalysisData{
		SourceUrl:     url,
		HtmlVersion:   "5",
		HeadingsCount: map[string]int{},
		LinkStats:     LinkStats{},
	}
}
