package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/service"
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

}

func (sh *ServiceHandler) ServiceStatus(w http.ResponseWriter, r *http.Request) {

}
