package analyzer

const (
	LoginForm = iota
	Unknown   = iota
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
	allLinks      map[string]bool
	LinkStats     LinkStats
	PageType      uint
}

type parsingState struct {
	currTag string
}

func (a *AnalysisData) ID() string {
	return a.SourceUrl
}

func NewAnalysis(url string) *AnalysisData {
	return &AnalysisData{
		SourceUrl:     url,
		HtmlVersion:   "5",
		allLinks:      map[string]bool{},
		HeadingsCount: map[string]int{},
		LinkStats:     LinkStats{},
	}
}
