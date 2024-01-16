package utils

import "net/url"

// GetAbsoluteURL converts relative or malformed URLs to absolute URLs
func GetAbsoluteURL(base *url.URL, href string) string {
	relURL, err := url.Parse(href)
	if err != nil {
		return ""
	}
	absURL := base.ResolveReference(relURL)
	return absURL.String()
}

// IsSameSite checks if two URLs belong to the same site
func IsSameSite(base *url.URL, target string) bool {
	targetURL, err := url.Parse(target)
	if err != nil {
		return false
	}
	return base.Host == targetURL.Host
}
