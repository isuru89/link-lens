package analyzer

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeUrlRelativeUrl(t *testing.T) {
	defer gock.Off()

	gock.New("https://www.google.com").
		Path("/path1").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString("<!doctype html><html><body><a href=\"/isuru\">rel link</a></body></html>")

	t.Run("Mock Test", func(t *testing.T) {
		info, err := AnalyzeUrl("https://www.google.com/path1")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.google.com/path1",
			HtmlVersion:   "5",
			Title:         "",
			HeadingsCount: map[string]int{},
			LinkStats: LinkStats{
				InternalLinks: 1,
				ExternalLinks: 0,
				InvalidLinks:  1,
			}},
			info)

	})

}

func TestAnalyzeUrl(t *testing.T) {

	// t.Run("HTML5 Test", func(t *testing.T) {
	// 	info, err := AnalyzeUrl("https://go.dev")
	// 	if err != nil {
	// 		t.Fatalf("Not suppose to throw an error! %v", err)
	// 	}

	// 	assert.Equal(t, &AnalysisData{
	// 		SourceUrl:   "https://go.dev",
	// 		HtmlVersion: "5",
	// 		Title:       "The Go Programming Language",
	// 		HeadingsCount: map[string]int{
	// 			"H1": 1, "H2": 4, "H3": 4,
	// 		},
	// 		LinkStats: LinkStats{
	// 			InternalLinks: 44,
	// 			ExternalLinks: 48,
	// 			InvalidLinks:  4,
	// 		}},
	// 		info)

	// })

	// t.Run("HTML4 Test", func(t *testing.T) {
	// 	info, err := AnalyzeUrl("https://www.w3.org/TR/1998/REC-CSS2-19980512/sample.html")
	// 	if err != nil {
	// 		t.Fatalf("Not suppose to throw an error! %v", err)
	// 	}

	// 	assert.Equal(t, &AnalysisData{
	// 		SourceUrl:     "https://www.w3.org/TR/1998/REC-CSS2-19980512/sample.html",
	// 		HtmlVersion:   "4",
	// 		Title:         "Appendix A: A sample style sheet for HTML 4.0",
	// 		HeadingsCount: map[string]int{"H1": 1},
	// 		LinkStats:     LinkStats{InternalLinks: 7, ExternalLinks: 0}},
	// 		info)

	// })

	// t.Run("JSON Endpoint Test", func(t *testing.T) {
	// 	_, err := AnalyzeUrl("https://jsonplaceholder.typicode.com/todos/1")
	// 	if err == nil {
	// 		t.Fatal("Suppose to throw an error!")
	// 	}

	// 	assert.Equal(t, "Only HTML content types are supported!", err.Error())
	// })
}
