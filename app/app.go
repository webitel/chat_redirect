package app

import (
	"chat_redirect/cache"
	"chat_redirect/model"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/webitel/wlog"
	"net/http"
)

type Application struct {
	config  *model.AppConfig
	handler *ApiHandler
	router  *mux.Router
	Redis   cache.CacheStore
}

func NewApplication(config *model.AppConfig) (*Application, error) {
	router := mux.NewRouter()
	redis, err := cache.NewRedisCache(*config.Redis.Host, *config.Redis.Port, *config.Redis.Password, *config.Redis.Database)
	if err != nil {
		return nil, err
	}
	app := &Application{
		config: config,
		router: router,
		Redis:  redis,
	}
	hndlr, err := NewApiHandler(app)
	if err != nil {
		return nil, err
	}
	app.handler = hndlr
	return app, nil
}

func (a *Application) Start() error {
	a.router.HandleFunc("/chat/redirect", a.handler.HandleRedirect)
	address := fmt.Sprintf("127.0.0.1:%d", *a.config.Server.Port)
	wlog.Debug(address)
	err := http.ListenAndServe(address, a.router)
	if err != nil {
		return err
	}
	return nil
}
