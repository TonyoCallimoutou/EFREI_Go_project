package domain

import "time"

type Shortener struct {
	ID        string    `json:"id" db:"id"`
	Url       string    `json:"url" db:"url"`
	ShortUrl  string    `json:"shortUrl" db:"shortUrl"`
	ExpiredAt time.Time `json:"expired_at" db:"expiredAt"`
	Count     int       `json:"count" db:"count"`
}
