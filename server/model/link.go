package model

import (
	_ "github.com/lib/pq"
)

type Link struct {
	ID          int    `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}
