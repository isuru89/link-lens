package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeUrl(t *testing.T) {
	t.Run("HTML5 test", func(t *testing.T) {
		info, err := AnalyzeUrl("https://go.dev")
		if err != nil {
			t.Fatalf("Not suppose to throw an error! %v", err)
		}

		assert.Equal(t, &AnalysisData{HtmlVersion: "5"}, info)

	})
}
