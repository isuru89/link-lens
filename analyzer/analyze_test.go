package analyzer

import (
	"errors"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeUrl_LinkTypesCrawl(t *testing.T) {
	defer gock.Off()

	// GIVEN
	gock.New("https://www.linklens.com").
		Path("/a/b/c").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html>
			<html>
			<title>Test Title</title>
			<body>
				<a href="/siterelative">site rel link</a>
				<a href="pathrelative/page1">path rel link</a>
				<a href="#anchor">anchor link</a>
				<a href="https://www.othersite.com/test/x">anchor link</a>
				<img src="https://www.amazons3.com/s3/image.png" alt="image" />
			</body>
			</html>`)

	gock.New("https://www.linklens.com").
		Path("/siterelative").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html><html></html>`)

	gock.New("https://www.linklens.com").
		Path("/a/b/pathrelative/page1").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html><html>page2</html>`)

	gock.New("https://www.othersite.com/test/x").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html><html>other-site</html>`)

	t.Run("Link Types Test", func(t *testing.T) {
		// WHEN
		info, err := AnalyzeUrl("https://www.linklens.com/a/b/c")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/a/b/c",
			HtmlVersion:   "5",
			Title:         "Test Title",
			HeadingsCount: map[string]int{},
			LinkStats: LinkStats{
				InternalLinks: 3,
				ExternalLinks: 1,
				InvalidLinks:  0,
			}},
			info)

	})
}

func TestAnalyzeUrl_InAccessibleLinks(t *testing.T) {
	defer gock.Off()

	// GIVEN
	gock.New("https://www.linklens.com").
		Path("/check/nx").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html>
			<html>
			<title>Test NX Links</title>
			<body>
				<a href="/siterelative">site rel link</a>
				<a href="/siterelative/nx">site rel link</a>
				<a href="pathrelative/page1">path rel link</a>
				<a href="pathrelative/pagenx">path rel link</a>
				<a href="pathrelative/pageerr">path error link</a>
				<a href="#anchor">anchor link</a>
				<a href="#anchor-nx">anchor link</a>
				<a href="https://www.othersite.com/test/x">anchor link</a>
				<a href="https://www.othersite.com/test/y/nx">anchor link</a>
			</body>
			</html>`)

	gock.New("https://www.linklens.com").
		Path("/siterelative").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html><html></html>`)

	gock.New("https://www.linklens.com").
		Path("/check/pathrelative/page1").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html><html>page2</html>`)

	gock.New("https://www.linklens.com").
		Path("/check/pathrelative/pagenx").
		Reply(404).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html><html>page 404</html>`)

	gock.New("https://www.linklens.com").
		Path("/check/pathrelative/pageerr").
		ReplyError(errors.New("Throwing error when page load"))

	gock.New("https://www.othersite.com/test/x").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html><html>other-site</html>`)

	t.Run("Link Inaccessibility Test", func(t *testing.T) {
		// WHEN
		info, err := AnalyzeUrl("https://www.linklens.com/check/nx")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/check/nx",
			HtmlVersion:   "5",
			Title:         "Test NX Links",
			HeadingsCount: map[string]int{},
			LinkStats: LinkStats{
				InternalLinks: 7,
				ExternalLinks: 2,
				InvalidLinks:  4, // anchor links will not be counted
			}},
			info)

	})
}

func TestAnalyzeUrl_HeadingCounts(t *testing.T) {
	defer gock.Off()

	// GIVEN
	gock.New("https://www.linklens.com").
		Path("/test/headings").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html>
			<html>
			<title>Test Headings</title>
			<body>
				<h1 /><h1>Heading 1</h1>
				<h2/><h2>Heading 2</h2>
				<h3 /><H3></H3>
				<h4/><h4></h4>
				<h5 /><h5>  </h5>
				<h6 /><H6></H6>
				<hr/><hr></hr>
			</body>
			</html>`)

	t.Run("Heading Count Test", func(t *testing.T) {
		// WHEN
		info, err := AnalyzeUrl("https://www.linklens.com/test/headings")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/headings",
			HtmlVersion:   "5",
			Title:         "Test Headings",
			HeadingsCount: map[string]int{"H1": 2, "H2": 2, "H3": 2, "H4": 2, "H5": 2, "H6": 2},
		}, info)

	})
}

func TestAnalyzeUrl_HtmlVersion(t *testing.T) {
	defer gock.Off()

	// GIVEN
	gock.New("https://www.linklens.com").
		Path("/test/htmlv4").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN">
			<html>
			<title>Test HTML V4</title>
			<body>
			</body>
			</html>`)

	gock.New("https://www.linklens.com").
		Path("/test/htmlv5").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!DOCTYPE HTML">
			<html>
			<title>Test HTML V5</title>
			<body>
			</body>
			</html>`)

	gock.New("https://www.linklens.com").
		Path("/test/htmlnx").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<html>
				<title>Test HTML Default</title>
				<body>
				</body>
				</html>`)

	t.Run("HTML V4 Test", func(t *testing.T) {
		// WHEN
		info, err := AnalyzeUrl("https://www.linklens.com/test/htmlv4")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/htmlv4",
			HtmlVersion:   "4",
			Title:         "Test HTML V4",
			HeadingsCount: map[string]int{},
		}, info)
	})

	t.Run("HTML V5 Test", func(t *testing.T) {
		// WHEN
		info, err := AnalyzeUrl("https://www.linklens.com/test/htmlv5")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/htmlv5",
			HtmlVersion:   "5",
			Title:         "Test HTML V5",
			HeadingsCount: map[string]int{},
		}, info)
	})

	t.Run("HTML V5 Default Test", func(t *testing.T) {
		// WHEN
		info, err := AnalyzeUrl("https://www.linklens.com/test/htmlnx")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/htmlnx",
			HtmlVersion:   "5",
			Title:         "Test HTML Default",
			HeadingsCount: map[string]int{},
		}, info)
	})
}

func TestAnalyzeUrl_TitleCapture(t *testing.T) {
	defer gock.Off()

	// GIVEN
	gock.New("https://www.linklens.com").
		Path("/test/title").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html>
			<html>
			<title>Test Title With</title>
			<body>
			</body>
			</html>`)

	gock.New("https://www.linklens.com").
		Path("/test/titlenx").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html>
				<html>
				<body>
				</body>
				</html>`)

	t.Run("Title Exists", func(t *testing.T) {
		// WHEN
		info, err := AnalyzeUrl("https://www.linklens.com/test/title")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/title",
			HtmlVersion:   "5",
			Title:         "Test Title With",
			HeadingsCount: map[string]int{},
		}, info)
	})

	t.Run("Title Does Not Exists", func(t *testing.T) {
		// WHEN
		info, err := AnalyzeUrl("https://www.linklens.com/test/titlenx")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/titlenx",
			HtmlVersion:   "5",
			Title:         "",
			HeadingsCount: map[string]int{},
		}, info)
	})
}
