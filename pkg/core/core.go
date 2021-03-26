package core

import (
	"net/url"
)

// IsValideAbsoluteURL checks if the URL provided is absolute and uses HTTP scheme.
func IsValideAbsoluteURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return false
	}

	return true
}
