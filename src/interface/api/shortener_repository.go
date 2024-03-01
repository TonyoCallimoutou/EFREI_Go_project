package api

import (
	"go_shortener/src/domain"
	"go_shortener/src/interface/database"
	"go_shortener/src/usecase"
)

type ShortenerController struct {
	Interactor usecase.ShortenerInteractor
}

func NewShortenerController(linkStore database.LinkStore) *ShortenerController {
	return &ShortenerController{
		Interactor: usecase.ShortenerInteractor{
			ShortenerRepository: &database.ShortenerRepository{
				LinkStore: linkStore,
			},
		},
	}
}

func (controller *ShortenerController) Create(url domain.Shortener) error {
	return controller.Interactor.Create(url)
}

func (controller *ShortenerController) GetById(url domain.Shortener) (string, error) {
	return controller.Interactor.GetById(url)
}

func (controller *ShortenerController) GetAll() (string, error) {
	return controller.Interactor.GetAll()
}

func (controller *ShortenerController) Update(url domain.Shortener) error {
	return controller.Interactor.Update(url)
}

func (controller *ShortenerController) Delete(url domain.Shortener) error {
	return controller.Interactor.Delete(url)
}
