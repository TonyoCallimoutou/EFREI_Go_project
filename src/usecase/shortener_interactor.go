package usecase

import "go_shortener/src/domain"

type ShortenerInteractor struct {
	ShortenerRepository ShortenerRepository
}

func (interactor *ShortenerInteractor) Create(shortener domain.Shortener) error {
	return interactor.ShortenerRepository.Create(shortener)
}

func (interactor *ShortenerInteractor) GetById(shortener domain.Shortener) (string, error) {
	return interactor.ShortenerRepository.GetById(shortener)
}

func (interactor *ShortenerInteractor) GetAll() (string, error) {
	return interactor.ShortenerRepository.GetAll()
}

func (interactor *ShortenerInteractor) Update(shortener domain.Shortener) error {
	return interactor.ShortenerRepository.Update(shortener)
}

func (interactor *ShortenerInteractor) Delete(shortener domain.Shortener) error {
	return interactor.ShortenerRepository.Delete(shortener)
}
