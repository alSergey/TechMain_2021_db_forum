package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/service"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type ServiceHandler struct {
	serviceUsecase service.ServiceUsecase
}

func NewServiceHandler(serviceUsecase service.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{
		serviceUsecase: serviceUsecase,
	}
}

func (sh *ServiceHandler) Configure(r *mux.Router) {
	r.HandleFunc("/service/clear", sh.ServiceClear).Methods(http.MethodPost)
	r.HandleFunc("/service/status", sh.ServiceStatus).Methods(http.MethodGet)
}

func (sh *ServiceHandler) ServiceClear(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("ServiceClear")
	errE := sh.serviceUsecase.ClearService()
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (sh *ServiceHandler) ServiceStatus(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("ServiceStatus")
	status, errE := sh.serviceUsecase.GetServiceStatus()
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, status, w)
}
