package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/user"
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

	user := &models.User{NickName: nickname}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	defer r.Body.Close()

	users, errE := uh.userUsecase.CreateUser(user)
	if errE != nil {
		errors.JSONSuccess(errE.HttpError, users, w)
		return
	}

	errors.JSONSuccess(http.StatusCreated, user, w)
}

func (uh *UserHandler) UserProfileGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user, errE := uh.userUsecase.GetUserByNickName(nickname)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, user, w)
}

func (uh *UserHandler) UserProfilePOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user := &models.User{NickName: nickname}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	defer r.Body.Close()

	errE := uh.userUsecase.EditUser(user)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, user, w)
}
