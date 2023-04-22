package domain

type App struct {
	AppName string `json:"app,omitempty"`
}

type AppUsecase interface {
	GetApp() *App
}
