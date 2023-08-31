package api

import (
	"github.com/AnSunn/ServerUserSegmentation/storage"
	_ "github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Trying to configure API instance (namely logger field)
func (a *API) configureLoggerField() error {
	log_level, err := logrus.ParseLevel(a.config.LogLevel)
	if err != nil {
		return err
	}
	a.logger.SetLevel(log_level)
	return nil
}

// Trying to configure router (namely router field)
func (a *API) configureRouterField() {
	a.router.HandleFunc(prefix+"/segments", a.PostSegment).Methods("POST")
	a.router.HandleFunc(prefix+"/segments", a.GetAllSegments).Methods("GET")
	a.router.HandleFunc(prefix+"/segments"+"/{title}", a.GetSegmentByTitle).Methods("GET")
	a.router.HandleFunc(prefix+"/segments"+"/{title}", a.DeleteSegmentByTitle).Methods("DELETE")
	a.router.HandleFunc(prefix+"/actions", a.PostAction).Methods("POST")
	a.router.HandleFunc(prefix+"/activesegments"+"/{user_id}", a.GetActiveSegmentsForUser).Methods("GET")
	a.router.HandleFunc(prefix+"/data"+"/{year}"+"/{month}", a.GetActionDataForAPeriod).Methods("GET")
}

// Trying to configure Storage (namely storage field)
func (a *API) configureStorageField() error {
	storage := storage.New(a.config.Store)
	if err := storage.Open(); err != nil {
		return err
	}
	a.storage = storage
	return nil
}
