package analyzer

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type AnalysisData struct {
	HtmlVersion  string
	Title        string
	HeadingsDist map[string]int
}

type parsingState struct {
	currTag string
}

var headingRegex = regexp.MustCompile(`(?i)h\d`)

func AnalyzeUrl(getUrl string) (*AnalysisData, error) {
	parsedUrl, err := url.Parse(getUrl)
	if err != nil {
		return nil, errors.New("Given URL is malformed!")
	}

	info := &AnalysisData{HtmlVersion: "5", HeadingsDist: map[string]int{}}
	errp := fetchUrlContent(parsedUrl, info)
	if errp != nil {
		return nil, errp
	}

	return info, nil
}

func fetchUrlContent(url *url.URL, info *AnalysisData) error {
	resp, err := http.Get(url.String())
	if err != nil {
		return errors.New("Cannt fetch the content from url!")
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return errors.New("Unsuccessful status code returned for the given URL!")
	}

	contentTypeHeader := resp.Header.Get("content-type")
	if strings.Index(strings.ToLower(contentTypeHeader), "text/html") < 0 {
		return errors.New("Only HTML content types are supported!")
	}

	t := html.NewTokenizer(resp.Body)
	status := &parsingState{}

	for {
		tokenType := t.Next()

		if tokenType == html.ErrorToken {
			if t.Err() == io.EOF {
				return nil
			}
			return t.Err()
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

		if tokenType == html.StartTagToken || tokenType == html.EndTagToken {
			node := t.Token()
			processToken(&node, info, status)
		}
	}
}

func processToken(token *html.Token, info *AnalysisData, status *parsingState) {
	if token.Type == html.StartTagToken {
		if token.Data == "title" {
			status.currTag = "title"
		} else if headingRegex.MatchString(token.Data) {
			info.HeadingsDist[strings.ToUpper(token.Data)]++
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
