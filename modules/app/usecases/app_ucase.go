package usecase

import "playground-go-api/domain"

type appUsecase struct{}

func NewAppUsecase() domain.AppUsecase {
	return &appUsecase{}
}

func (a *appUsecase) GetApp() *domain.App {
	return &domain.App{
		AppName: "playground-go-api",
	}
}
