package http

import (
	"encoding/json"
	"fmt"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/post"
)

type PostHandler struct {
	postUsecase post.PostUsecase
}

func NewPostHandler(postUsecase post.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

func (ph *PostHandler) Configure(r *mux.Router) {
	r.HandleFunc("/post/{id}/details", ph.PostDetailsGET).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}/details", ph.PostDetailsPOST).Methods(http.MethodPost)
}

func (ph *PostHandler) PostDetailsGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostDetailsGET")
	vars := mux.Vars(r)
	strId := vars["id"]
	//fmt.Println("PostDetailsGET strId = ", strId)

	id, err := strconv.Atoi(strId)
	if err != nil {
		return
	}
	//fmt.Println("PostDetailsGET id = ", id)

	related := r.URL.Query().Get("related")
	fullPostParams := &models.FullPostParams{}

	if strings.Contains(related, "user") {
		fullPostParams.User = true
	}

	if strings.Contains(related, "forum") {
		fullPostParams.Forum = true
	}

	if strings.Contains(related, "thread") {
		fullPostParams.Thread = true
	}
	//fmt.Println("PostDetailsGET params = ", fullPostParams)

	post, errE := ph.postUsecase.GetPost(id, fullPostParams)
	if errE != nil {
		//fmt.Println("PostDetailsGET error = ", errE)
		errors.JSONError(errE, w)
		return
	}

	//fmt.Println("PostDetailsGET result post = ", post)
	errors.JSONSuccess(http.StatusOK, post, w)
}

func (ph *PostHandler) PostDetailsPOST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostDetailsPOST")
	vars := mux.Vars(r)
	strId := vars["id"]
	//fmt.Println("PostDetailsPOST strId = ", strId)

	id, err := strconv.Atoi(strId)
	if err != nil {
		return
	}
	//fmt.Println("PostDetailsPOST id = ", id)

	post := &models.Post{Id: id}
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		return
	}
	defer r.Body.Close()
	//fmt.Println("PostDetailsPOST post = ", post)

	post, errE := ph.postUsecase.UpdatePost(post)
	if errE != nil {
		//fmt.Println("PostDetailsPOST error = ", errE)
		errors.JSONError(errE, w)
		return
	}

	//fmt.Println("PostDetailsPOST result post = ", post)
	errors.JSONSuccess(http.StatusOK, post, w)
}
