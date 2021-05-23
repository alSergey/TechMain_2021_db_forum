package http

import (
	"encoding/json"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user/model"
)

type UserHandler struct {
	userUsecase user.UserUsecase
}

func NewUserHandler(sserUsecase user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: sserUsecase,
	}
}

func (uh *UserHandler) Configure(r *mux.Router) {
	r.HandleFunc("/user/{nickname}/create", uh.UserCreate).Methods(http.MethodPost)
	r.HandleFunc("/user/{nickname}/profile", uh.UserProfileGET).Methods(http.MethodGet)
	r.HandleFunc("/user/{nickname}/profile", uh.UserProfilePOST).Methods(http.MethodPost)
}

func (uh *UserHandler) UserCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	user.NickName = nickname

	users, errE := uh.userUsecase.Create(user)
	if errE != nil {
		errors.JSONSuccess(errE.HttpError, users, w)
		return
	}

	errors.JSONSuccess(http.StatusCreated, user, w)
}

func (uh *UserHandler) UserProfileGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user, errE := uh.userUsecase.GetByNickName(nickname)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, user, w)
}

func (uh *UserHandler) UserProfilePOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	user.NickName = nickname

	errE := uh.userUsecase.Edit(user)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, user, w)
}
