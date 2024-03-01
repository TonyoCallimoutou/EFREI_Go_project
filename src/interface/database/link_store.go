package database

import "go_shortener/src/domain"

type LinkStore interface {
	Create(url domain.Shortener) error
	GetById(url domain.Shortener) (string, error)
	GetAll() (string, error)
	Update(url domain.Shortener) error
	Delete(url domain.Shortener) error
}
