package analyzer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var headingRegex = regexp.MustCompile(`(?i)h\d`)

func AnalyzeUrl(getUrl string) (*AnalysisData, error) {
	log.Printf("Starting the anlysis of url: %s", getUrl)

	parsedUrl, err := url.Parse(getUrl)
	if err != nil {
		return nil, errors.New("Given URL is malformed!")
	}

	info := NewAnalysis(getUrl)
	errp := fetchUrlContent(parsedUrl, info)
	if errp != nil {
		return nil, errp
	}

	Crawl(info)

	info.PageType = Unknown
	info.allLinks = nil

	return info, nil
}

func fetchUrlContent(url *url.URL, info *AnalysisData) error {
	resp, err := http.Get(url.String())
	if err != nil {
		return errors.New("Cannot fetch the content from url!")
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Unsuccessful status code returned for the given URL! %d", resp.StatusCode))
	}

	contentTypeHeader := resp.Header.Get("content-type")
	if strings.Index(strings.ToLower(contentTypeHeader), "text/html") < 0 {
		return errors.New("Only HTML content types are supported!")
	}

	log.Printf("Recieved a valid html content from %s", url.String())
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
				if strings.ToLower(v.Key) == "href" {
					info.allLinks[v.Val] = true
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
