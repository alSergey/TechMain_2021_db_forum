package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/forum"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
)

type ForumHandler struct {
	forumUsecase forum.ForumUsecase
}

func NewForumHandler(forumUsecase forum.ForumUsecase) *ForumHandler {
	return &ForumHandler{
		forumUsecase: forumUsecase,
	}
}

func (fh *ForumHandler) Configure(r *mux.Router) {
	r.HandleFunc("/forum/create", fh.ForumCreate).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/details", fh.ForumDetails).Methods(http.MethodGet)
	r.HandleFunc("/forum/{slug}/create", fh.ForumCreateThread).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/users", fh.ForumUsers).Methods(http.MethodGet)
	r.HandleFunc("/forum/{slug}/threads", fh.ForumThreads).Methods(http.MethodGet)
}

func (fh *ForumHandler) ForumCreate(w http.ResponseWriter, r *http.Request) {
	forum := &models.Forum{}
	err := json.NewDecoder(r.Body).Decode(&forum)
	if err != nil {
		return
	}
	defer r.Body.Close()

	existForum, errE := fh.forumUsecase.Create(forum)
	if errE != nil {
		switch errE.ErrorCode {
		case errors.ForumCreateNotExist:
			errors.JSONError(errE, w)
			return

		case errors.ForumCreateConflict:
			errors.JSONSuccess(errE.HttpError, existForum, w)
			return
		}
	}

	errors.JSONSuccess(http.StatusCreated, forum, w)
}

func (fh *ForumHandler) ForumDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	forum, errE := fh.forumUsecase.GetBySlug(slug)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, forum, w)
}

func (fh *ForumHandler) ForumCreateThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	thread := &models.Thread{Forum: slug}
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		return
	}
	defer r.Body.Close()

	isSlug := thread.Slug != ""

	existThread, errE := fh.forumUsecase.CreateThread(thread)
	if errE != nil {
		switch errE.ErrorCode {
		case errors.ForumCreateThreadNotExist:
			errors.JSONError(errE, w)
			return

		case errors.ForumCreateThreadConflict:
			errors.JSONSuccess(errE.HttpError, existThread, w)
			return
		}
	}

	if !isSlug {
		errors.JSONSuccess(http.StatusCreated, models.ConvertThread(thread), w)
		return
	}

	errors.JSONSuccess(http.StatusCreated, thread, w)
}

func (fh *ForumHandler) ForumUsers(w http.ResponseWriter, r *http.Request) {

}

func (fh *ForumHandler) ForumThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	threadsParams := &models.ThreadParams{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(threadsParams, r.URL.Query())
	if err != nil {
		return
	}

	threads, errE := fh.forumUsecase.GetThreadsBySlugAndParams(slug, threadsParams)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, threads, w)
}
