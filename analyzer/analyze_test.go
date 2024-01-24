package analyzer

import (
	"errors"
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeUrl_Errors(t *testing.T) {
	defer gock.Off()

	tests := map[string]struct {
		preRun    func()
		url       string
		errorCode string
		errorMsg  string
	}{
		"Invalid URL": {
			preRun:    func() {},
			url:       ":invalid url",
			errorCode: ErrorInvalidUrl,
			errorMsg:  "given url is malformed",
		},
		"Non Existence URL": {
			preRun:    func() {},
			url:       "http://non-exist.ne",
			errorCode: RemoteFetchError,
			errorMsg:  "cannot fetch the content from url",
		},
		"Returns != 200": {
			preRun: func() {
				gock.New("https://www.othersite.com/test/x").
					Reply(404).
					AddHeader("content-type", "text/html").
					BodyString(`<!doctype html><html>404 Not Exists</html>`)
			},
			url:       "https://www.othersite.com/test/x",
			errorCode: UnsuccessfulStatusCode,
			errorMsg:  "unsuccessful status code returned for the given url! 404",
		},
		"Non HTML": {
			preRun: func() {
				gock.New("https://www.othersite.com/test/y").
					Reply(200).
					AddHeader("content-type", "application/json").
					BodyString(`{ "alive": true }`)
			},
			url:       "https://www.othersite.com/test/y",
			errorCode: InvalidContentType,
			errorMsg:  "only HTML content types are supported",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			test.preRun()

			// WHEN
			_, err := AnalyzeUrl(test.url, &OneDepthCrawler{})

			// THEN
			if err == nil {
				assert.Fail(t, "Expected to fail, but did not")
			}
			assert.Equal(t, fmt.Sprintf("[%s] %s", test.errorCode, test.errorMsg), err.Error())
		})
	}
}

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
	gock.New("https://www.othersite.com/test/x").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html><html>other-site</html>`)

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
				InternalLinkCount: 3,
				ExternalLinkCount: 1,
				InvalidLinkCount:  0,
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
			<a href="/st-relative/nx">site rel link</a>
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
	mockHtmlUrlWithStatusCode("/check/pathrelative/pagenx", `<!doctype html><html>page 404</html>`, 404)
	gock.New("https://www.linklens.com").
		Path("/check/pathrelative/pageerr").
		ReplyError(errors.New("Throwing error when page load"))
	gock.New("https://www.othersite.com/test/x").
		Reply(200).
		AddHeader("content-type", "text/html").
		BodyString(`<!doctype html><html>External Site</html>`)

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
				InternalLinkCount: 7,
				ExternalLinkCount: 2,
				InvalidLinkCount:  4, // anchor links will not be counted
				InvalidLinks: []string{
					"https://www.linklens.com/check/pathrelative/pageerr",
					"https://www.linklens.com/check/pathrelative/pagenx",
					"https://www.linklens.com/st-relative/nx",
					"https://www.othersite.com/test/y/nx",
				},
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

	tests := map[string]struct {
		preRun              func()
		url                 string
		expectedHtmlVersion string
	}{
		"HTML V4 Test": {
			preRun: func() {
				mockHtmlUrl("/test/htmlv4", `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN">
				<html><title>Test HTML V4</title><body></body></html>`)
			},
			url:                 "https://www.linklens.com/test/htmlv4",
			expectedHtmlVersion: "4",
		},
		"HTML V5 Test": {
			preRun: func() {
				mockHtmlUrl("/test/htmlv5", `<!DOCTYPE HTML"><html><title>Test HTML V5</title><body></body></html>`)
			},
			url:                 "https://www.linklens.com/test/htmlv5",
			expectedHtmlVersion: "5",
		},
		"Default Should Be V5": {
			preRun: func() {
				mockHtmlUrl("/test/htmlnx", `<html><title>Test HTML Default</title><body></body></html>`)
			},
			url:                 "https://www.linklens.com/test/htmlnx",
			expectedHtmlVersion: "5",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			test.preRun()

			// WHEN
			info := callAnalysisUrlSuccess(t, test.url)

			// THEN
			assert.Equal(t, test.expectedHtmlVersion, info.HtmlVersion)
		})
	}
}

func TestAnalyzeUrl_TitleCapture(t *testing.T) {
	defer gock.Off()

	tests := map[string]struct {
		preRun        func()
		url           string
		expectedTitle string
	}{
		"Title Exists": {
			preRun: func() {
				mockHtmlUrl("/test/title", `<!doctype html><html><title>Test Title With</title><body></body></html>`)
			},
			url:           "https://www.linklens.com/test/title",
			expectedTitle: "Test Title With",
		},
		"Title Does Not Exists": {
			preRun: func() {
				mockHtmlUrl("/test/titlenx", `<!doctype html><html><body></body></html>`)
			},
			url:           "https://www.linklens.com/test/titlenx",
			expectedTitle: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			test.preRun()

			// WHEN
			info := callAnalysisUrlSuccess(t, test.url)

			// THEN
			assert.Equal(t, test.expectedTitle, info.Title)
		})
	}
}

func TestAnalyzeUrl_IsLoginForm(t *testing.T) {
	defer gock.Off()

	// GIVEN
	mockHtmlUrl("/test/loginform", `<html>
		<title>Test Login Form</title>
		<body>
			<form>
				<input type="text" name="email"></input>
				<INPUT TYPE="password" name="password"></input>
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

	expected := func(url string, pageType string) *AnalysisData {
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
		pageType string
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
	info, err := AnalyzeUrl(url, &OneDepthCrawler{})
	if err != nil {
		t.Fatalf("Not suppose to throw an error! %v", err)
	}
	return info
}
