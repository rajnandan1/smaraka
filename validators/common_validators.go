package validators

import (
	_url "net/url"
)

// check if valid url or not
func IsValidURL(url string) bool {
	_, err := _url.ParseRequestURI(url)
	if err != nil {
		return false
	}
	return true
}

// write a strong password validator
func IsValidPassword(password string) bool {
	// Check minimum length
	if len(password) < 8 {
		return false
	}

	// Return true only if all requirements are met
	return true
}
