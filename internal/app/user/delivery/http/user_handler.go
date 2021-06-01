package http

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("UserCreate")
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	//fmt.Println("UserCreate nickname = ", nickname)

	user := &models.User{NickName: nickname}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	defer r.Body.Close()
	//fmt.Println("UserCreate user = ", user)

	existUser, errE := uh.userUsecase.CreateUser(user)
	if errE != nil {
		//fmt.Println("UserCreate exist user = ", existUser)
		errors.JSONSuccess(errE.HttpError, existUser, w)
		return
	}
	//fmt.Println("UserCreate result user = ", user)

	errors.JSONSuccess(http.StatusCreated, user, w)
}

func (uh *UserHandler) UserProfileGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UserProfileGET")
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	//fmt.Println("UserProfileGET nickname = ", nickname)

	user, errE := uh.userUsecase.GetUserByNickName(nickname)
	if errE != nil {
		//fmt.Println("UserProfileGET error = ", errE)
		errors.JSONError(errE, w)
		return
	}
	//fmt.Println("UserProfileGET result user = ", user)

	errors.JSONSuccess(http.StatusOK, user, w)
}

func (uh *UserHandler) UserProfilePOST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UserProfilePOST")
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	//fmt.Println("UserProfilePOST nickname = ", nickname)

	user := &models.User{NickName: nickname}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	defer r.Body.Close()
	//fmt.Println("UserProfilePOST user = ", user)

	errE := uh.userUsecase.EditUser(user)
	if errE != nil {
		//fmt.Println("UserProfilePOST error = ", errE)
		errors.JSONError(errE, w)
		return
	}
	//fmt.Println("UserProfilePOST result user = ", user)

	errors.JSONSuccess(http.StatusOK, user, w)
}
