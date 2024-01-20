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
	mockHtmlUrl("/a/b/c", `<!doctype html>
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
	mockHtmlUrl("/siterelative", `<!doctype html><html></html>`)
	mockHtmlUrl("/a/b/pathrelative/page1", `<!doctype html><html>page2</html>`)
	mockHtmlUrl("/test/x", `<!doctype html><html>other-site</html>`)

	t.Run("Link Types Test", func(t *testing.T) {
		// WHEN
		info := callAnalysisUrlSuccess(t, "https://www.linklens.com/a/b/c")

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
			},
			PageType: Unknown,
		}, info)
	})
}

func TestAnalyzeUrl_InAccessibleLinks(t *testing.T) {
	defer gock.Off()

	// GIVEN
	mockHtmlUrl("/check/nx", `<!doctype html>
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
	mockHtmlUrl("/siterelative", `<!doctype html><html></html>`)
	mockHtmlUrl("/check/pathrelative/page1", `<!doctype html><html>page2</html>`)
	mockHtmlUrl("/test/x", `<!doctype html><html>other-site</html>`)
	mockHtmlUrlWithStatusCode("/check/pathrelative/pagenx", `<!doctype html><html>page 404</html>`, 404)
	gock.New("https://www.linklens.com").
		Path("/check/pathrelative/pageerr").
		ReplyError(errors.New("Throwing error when page load"))

	t.Run("Link Inaccessibility Test", func(t *testing.T) {
		// WHEN
		info := callAnalysisUrlSuccess(t, "https://www.linklens.com/check/nx")

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
			},
			PageType: Unknown,
		}, info)
	})
}

func TestAnalyzeUrl_HeadingCounts(t *testing.T) {
	defer gock.Off()

	// GIVEN
	mockHtmlUrl("/test/headings", `<!doctype html>
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
		info := callAnalysisUrlSuccess(t, "https://www.linklens.com/test/headings")

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/headings",
			HtmlVersion:   "5",
			Title:         "Test Headings",
			HeadingsCount: map[string]int{"H1": 2, "H2": 2, "H3": 2, "H4": 2, "H5": 2, "H6": 2},
			PageType:      Unknown,
		}, info)
	})
}

func TestAnalyzeUrl_HtmlVersion(t *testing.T) {
	defer gock.Off()

	// GIVEN
	mockHtmlUrl("/test/htmlv4", `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN">
		<html><title>Test HTML V4</title><body></body></html>`)
	mockHtmlUrl("/test/htmlv5", `<!DOCTYPE HTML"><html><title>Test HTML V5</title><body></body></html>`)
	mockHtmlUrl("/test/htmlnx", `<html><title>Test HTML Default</title><body></body></html>`)

	t.Run("HTML V4 Test", func(t *testing.T) {
		// WHEN
		info := callAnalysisUrlSuccess(t, "https://www.linklens.com/test/htmlv4")

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/htmlv4",
			HtmlVersion:   "4",
			Title:         "Test HTML V4",
			HeadingsCount: map[string]int{},
			PageType:      Unknown,
		}, info)
	})

	t.Run("HTML V5 Test", func(t *testing.T) {
		// WHEN
		info := callAnalysisUrlSuccess(t, "https://www.linklens.com/test/htmlv5")

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/htmlv5",
			HtmlVersion:   "5",
			Title:         "Test HTML V5",
			HeadingsCount: map[string]int{},
			PageType:      Unknown,
		}, info)
	})

	t.Run("HTML V5 Default Test", func(t *testing.T) {
		// WHEN
		info := callAnalysisUrlSuccess(t, "https://www.linklens.com/test/htmlnx")

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/htmlnx",
			HtmlVersion:   "5",
			Title:         "Test HTML Default",
			HeadingsCount: map[string]int{},
			PageType:      Unknown,
		}, info)
	})
}

func TestAnalyzeUrl_TitleCapture(t *testing.T) {
	defer gock.Off()

	// GIVEN
	mockHtmlUrl("/test/title", `<!doctype html><html><title>Test Title With</title><body></body></html>`)
	mockHtmlUrl("/test/titlenx", `<!doctype html><html><body></body></html>`)

	t.Run("Title Exists", func(t *testing.T) {
		// WHEN
		info := callAnalysisUrlSuccess(t, "https://www.linklens.com/test/title")

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/title",
			HtmlVersion:   "5",
			Title:         "Test Title With",
			HeadingsCount: map[string]int{},
			PageType:      Unknown,
		}, info)
	})

	t.Run("Title Does Not Exists", func(t *testing.T) {
		// WHEN
		info := callAnalysisUrlSuccess(t, "https://www.linklens.com/test/titlenx")

		// THEN
		assert.Equal(t, &AnalysisData{
			SourceUrl:     "https://www.linklens.com/test/titlenx",
			HtmlVersion:   "5",
			Title:         "",
			HeadingsCount: map[string]int{},
			PageType:      Unknown,
		}, info)
	})
}

func TestAnalyzeUrl_IsLoginForm(t *testing.T) {
	defer gock.Off()

	// GIVEN
	mockHtmlUrl("/test/loginform", `<html>
		<title>Test Login Form</title>
		<body>
			<form>
				<input type="text" name="email"></input>
				<input type="password" name="password"></input>
				<input type="submit">Login</button>
			</form>
		</body>
		</html>`)
	mockHtmlUrl("/test/loginsubmit", `<html>
		<title>Test Login Form</title>
		<body>
			<form>
				<input type="text" name="email"></input>
				<input type="submit">Login</button>
			</form>
		</body>
		</html>`)
	mockHtmlUrl("/test/loginpw", `<html>
		<title>Test Login Form</title>
		<body>
			<form>
				<input type="text" name="email"></input>
				<input type="password" name="password"></input>
			</form>
		</body>
		</html>`)
	mockHtmlUrl("/test/loginmultiplepw", `<html>
		<title>Test Login Form</title>
		<body>
			<form>
				<input type="text" name="email"></input>
				<input type="password" name="password1"></input>
				<input type="password" name="password2"></input>
				<input type="submit">Login</button>
			</form>
		</body>
		</html>`)
	mockHtmlUrl("/test/loginmultiplesubmit", `<html>
		<title>Test Login Form</title>
		<body>
			<form>
				<input type="text" name="email"></input>
				<input type="password" name="password"></input>
				<input type="submit">Login 1</button>
				<input type="submit">Login 2</button>
			</form>
		</body>
		</html>`)

	expected := func(url string, pageType uint) *AnalysisData {
		return &AnalysisData{
			SourceUrl:     url,
			HtmlVersion:   "5",
			Title:         "Test Login Form",
			HeadingsCount: map[string]int{},
			PageType:      pageType,
		}
	}

	testcases := map[string]struct {
		url      string
		pageType uint
	}{
		"Should Be A Login Form":                  {pageType: LoginForm, url: "https://www.linklens.com/test/loginform"},
		"No Login Form: Only Submit Button":       {pageType: Unknown, url: "https://www.linklens.com/test/loginsubmit"},
		"No Login Form: Only Password Input":      {pageType: Unknown, url: "https://www.linklens.com/test/loginpw"},
		"No Login Form: Multiple Password Inputs": {pageType: Unknown, url: "https://www.linklens.com/test/loginmultiplepw"},
		"No Login Form: Multiple Submits":         {pageType: Unknown, url: "https://www.linklens.com/test/loginmultiplesubmit"},
	}

	for name, tcase := range testcases {
		t.Run(name, func(t *testing.T) {
			// WHEN
			info := callAnalysisUrlSuccess(t, tcase.url)

			// THEN
			assert.Equal(t, expected(tcase.url, tcase.pageType), info)
		})
	}
}

func mockHtmlUrl(path, response string) {
	mockHtmlUrlWithStatusCode(path, response, 200)
}

func mockHtmlUrlWithStatusCode(path, response string, statusCode int) {
	gock.New("https://www.linklens.com").
		Path(path).
		Reply(statusCode).
		AddHeader("content-type", "text/html").
		BodyString(response)
}

func callAnalysisUrlSuccess(t *testing.T, url string) *AnalysisData {
	info, err := AnalyzeUrl(url)
	if err != nil {
		t.Fatalf("Not suppose to throw an error! %v", err)
	}
	return info
}
