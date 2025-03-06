package models

import "time"

type Bookmark struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Excerpt     string    `json:"excerpt"`
	ImageSmall  string    `json:"image_small"`
	ImageLarge  string    `json:"image_large"`
	GroupID     string    `json:"group_id"`
	Status      string    `json:"status"`
	FullText    string    `json:"full_text"`
	AccentColor string    `json:"accent_color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GithubRepo struct {
	URL     string `json:"url"`
	Stars   int    `json:"stars"`
	Excerpt string `json:"excerpt"`
	Title   string `json:"title"`
}
