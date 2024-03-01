package domain

import "time"

type Shortener struct {
	ID        string    `json:"id"`
	Url       string    `json:"url"`
	ShortUrl  string    `json:"shortUrl"`
	ExpiredAt time.Time `json:"expired_at"`
}
