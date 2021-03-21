package core_test

import (
	"testing"

	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestExtractParentURL(t *testing.T) {
	tests := map[string]struct {
		rawURL         string
		expectedOutput bool
	}{
		"empty url":    {rawURL: "", expectedOutput: false},
		"absolute url": {rawURL: "http://feeds.bbci.co.uk/news/uk/rss.xml", expectedOutput: true},
		"relative url": {rawURL: "/path/to/file", expectedOutput: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value := core.IsValideAbsoluteURL(test.rawURL)
			assert.Equal(t, test.expectedOutput, value)
		})
	}
}
