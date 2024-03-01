package database

import "go_shortener/src/domain"

type ShortenerRepository struct {
	LinkStore
}

func (db *ShortenerRepository) Create(shortener domain.Shortener) error {
	return db.Create(shortener)
}

func (db *ShortenerRepository) getById() (string, error) {
	return db.getById()
}

func (db *ShortenerRepository) getAll() (string, error) {
	return db.GetAll()
}

func (db *ShortenerRepository) Update(shortener domain.Shortener) error {
	return db.Update(shortener)
}

func (db *ShortenerRepository) Delete(shortener domain.Shortener) error {
	return db.Delete(shortener)
}
