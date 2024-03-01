package usecase

import "go_shortener/src/domain"

type ShortenerRepository interface {
	Create(url domain.Shortener) error
	GetById(url domain.Shortener) (string, error)
	GetAll() (string, error)
	Update(url domain.Shortener) error
	Delete(url domain.Shortener) error
}
