package analyzer

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var headingRegex = regexp.MustCompile(`(?i)h\d`)

func AnalyzeUrl(getUrl string, crawler Crawler) (*AnalysisData, error) {
	slog.Info("Starting the anlysis of ", "url", getUrl, "crawler", reflect.TypeOf(crawler).Elem())

	parsedUrl, err := url.Parse(getUrl)
	if err != nil {
		return nil, &AnalysisError{
			ErrorCode: ErrorInvalidUrl,
			Cause:     fmt.Errorf("given url is malformed"),
		}
	}

	info := NewAnalysis(getUrl)
	status, errp := fetchUrlContent(parsedUrl, info)
	if errp != nil {
		return nil, errp
	}

	// crawl links
	stats := crawler.Crawl(info.SourceUrl, status.allLinks)
	info.LinkStats = *stats

	// guess page type...
	derivePageType(status, info)

	return info, nil
}

func fetchUrlContent(url *url.URL, info *AnalysisData) (*parsingState, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, &AnalysisError{
			ErrorCode: RemoteFetchError,
			Cause:     fmt.Errorf("cannot fetch the content from url"),
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, &AnalysisError{
			ErrorCode: UnsuccessfulStatusCode,
			Cause:     fmt.Errorf("unsuccessful status code returned for the given url! %d", resp.StatusCode),
		}
	}

	contentTypeHeader := resp.Header.Get("content-type")
	if !strings.Contains(strings.ToLower(contentTypeHeader), "text/html") {
		return nil, &AnalysisError{
			ErrorCode: InvalidContentType,
			Cause:     fmt.Errorf("only HTML content types are supported"),
		}
	}

	slog.Info("Recieved a valid html content from ", "url", url.String())
	t := html.NewTokenizer(resp.Body)
	status := &parsingState{inputTypeCounts: map[string]int{}, allLinks: map[string]bool{}}

	for {
		tokenType := t.Next()

		if tokenType == html.ErrorToken {
			if t.Err() == io.EOF {
				return status, nil
			}
			return status, t.Err()
		}

		if tokenType == html.DoctypeToken {
			html4Regex := regexp.MustCompile(`(?i)HTML\s+[4].*`)
			doctag := string(t.Text())

			result := html4Regex.FindStringSubmatch(doctag)
			if result != nil {
				info.HtmlVersion = "4"
			}
		}

		if tokenType == html.TextToken {
			processText(string(t.Text()), info, status)
		}

		if tokenType == html.StartTagToken || tokenType == html.EndTagToken || tokenType == html.SelfClosingTagToken {
			node := t.Token()
			processToken(&node, info, status)
		}
	}
}

func processToken(token *html.Token, info *AnalysisData, status *parsingState) {
	if token.Type == html.StartTagToken || token.Type == html.SelfClosingTagToken {
		if token.Data == "title" {
			status.currTag = "title"
		} else if headingRegex.MatchString(token.Data) {
			info.HeadingsCount[strings.ToUpper(token.Data)]++
		} else if token.Data == "a" {
			for _, v := range token.Attr {
				if v.Key == "href" {
					status.allLinks[v.Val] = true
					break
				}
			}
		} else if token.Data == "input" {
			for _, v := range token.Attr {
				if v.Key == "type" {
					status.inputTypeCounts[v.Val]++
					break
				}
			}
		}
	} else if token.Type == html.EndTagToken {
		if status.currTag != "" {
			status.currTag = ""
		}
	}
}

func processText(content string, info *AnalysisData, status *parsingState) {
	if status.currTag == "title" {
		info.Title = content
	}
}

func derivePageType(status *parsingState, info *AnalysisData) {
	if status.inputTypeCounts["password"] == 1 && status.inputTypeCounts["submit"] == 1 {
		info.PageType = LoginForm
	} else {
		info.PageType = Unknown
	}
}
