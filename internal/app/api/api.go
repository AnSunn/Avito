package api

import (
	"github.com/AnSunn/ServerUserSegmentation/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	prefix string = "/api/v1"
)

// Base API server instance description
type API struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	//add field to work with Storage
	storage *storage.Storage
}

// API constructor: build base API instance
func New(config *Config) *API {
	return &API{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start http server/configure loggers/router/db connection/etc
func (api *API) Start() error {
	//Trying to configure Logger
	if err := api.configureLoggerField(); err != nil {
		return err
	}
	//Confirmation that logger is successfully configured
	api.logger.Info("Starting API server at port", api.config.BindAddr)

	//Configure router
	api.configureRouterField()

	//Configure storage
	if err := api.configureStorageField(); err != nil {
		return err
	}
	//If validation is successfully over, the listenandserve has to be started
	return http.ListenAndServe(api.config.BindAddr, api.router)
}
