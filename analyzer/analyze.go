package analyzer

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"golang.org/x/net/html"
)

type AnalysisData struct {
	HtmlVersion string
}

func AnalyzeUrl(getUrl string) (*AnalysisData, error) {
	parsedUrl, err := url.Parse(getUrl)
	if err != nil {
		return nil, errors.New("Given URL is malformed!")
	}

	info := &AnalysisData{HtmlVersion: "5"}
	urlContent, err := fetchUrlContent(parsedUrl, info)
	if err != nil {
		return nil, err
	}

	fmt.Print(urlContent)
	return info, nil
}

func fetchUrlContent(url *url.URL, info *AnalysisData) (string, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return "", errors.New("Cannt fetch the content from url!")
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", errors.New("Unsuccessful status code returned for the given URL!")
	}

	t := html.NewTokenizer(resp.Body)

	for {
		tokenType := t.Next()

		if tokenType == html.ErrorToken {
			if t.Err() == io.EOF {
				return "", nil
			}
			return "", t.Err()
		}

		if tokenType == html.DoctypeToken {
			html4Regex := regexp.MustCompile(`(?i)HTML\s+[4].*`)
			doctag := string(t.Text())

			result := html4Regex.FindStringSubmatch(doctag)
			if result != nil {
				info.HtmlVersion = "4"
			}
		}

		if tokenType == html.StartTagToken {
			// node := t.Token()
			// fmt.Println(string(t.Raw()))
			// fmt.Println(string(t.Text()))
		}
	}
}
