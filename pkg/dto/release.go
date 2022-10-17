//nolint:tagliatelle
package dto

import "time"

// The Release structure contains only the necessary keys.
type Release struct {
	ID          int       `json:"id"`
	HTMLURL     string    `json:"html_url"`
	Name        string    `json:"name"`
	TagName     string    `json:"tag_name"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
}
