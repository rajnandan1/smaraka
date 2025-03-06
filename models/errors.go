package models

import "errors"

var (
	ErrorIndexPathDoesNotExist = errors.New("index is not initialized")
	ErrorNotOfTypeBookmark     = errors.New("document is not of type Bookmark")
)
