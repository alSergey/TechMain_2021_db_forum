package http

import (
	"encoding/json"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/models"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/post"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/thread"
)

type ThreadHandler struct {
	threadUsecase thread.ThreadUsecase
	postUsecase   post.PostUsecase
}

func NewThreadHandler(threadUsecase thread.ThreadUsecase, postUsecase post.PostUsecase) *ThreadHandler {
	return &ThreadHandler{
		threadUsecase: threadUsecase,
		postUsecase:   postUsecase,
	}
}

func (th *ThreadHandler) Configure(r *mux.Router) {
	r.HandleFunc("/thread/{slug_or_id}/create", th.ThreadCreate).Methods(http.MethodPost)
	r.HandleFunc("/thread/{slug_or_id}/details", th.ThreadDetailsGET).Methods(http.MethodGet)
	r.HandleFunc("/thread/{slug_or_id}/details", th.ThreadDetailsPOST).Methods(http.MethodPost)
	r.HandleFunc("/thread/{slug_or_id}/posts", th.ThreadPosts).Methods(http.MethodGet)
	r.HandleFunc("/thread/{slug_or_id}/vote", th.ThreadVote).Methods(http.MethodPost)
}

func (th *ThreadHandler) ThreadCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]

	var posts []*models.Post
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		return
	}
	defer r.Body.Close()

	resultPosts, errE := th.postUsecase.CreatePost(slug, posts)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusCreated, resultPosts, w)
}

func (th *ThreadHandler) ThreadDetailsGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]

	thread, errE := th.threadUsecase.GetThread(slug)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, thread, w)
}

func (th *ThreadHandler) ThreadDetailsPOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]

	thread := &models.Thread{Forum: slug}
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		return
	}
	defer r.Body.Close()

	thread, errE := th.threadUsecase.UpdateThread(thread)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, thread, w)
}

func (th *ThreadHandler) ThreadPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]

	postParams := &models.PostParams{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(postParams, r.URL.Query())
	if err != nil {
		return
	}

	posts, errE := th.postUsecase.GetPostsBySlugAndParams(slug, postParams)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, posts, w)
}

func (th *ThreadHandler) ThreadVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]

	vote := &models.Vote{}
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		return
	}
	defer r.Body.Close()

	thread, errE := th.threadUsecase.Vote(slug, vote)
	if errE != nil {
		errors.JSONError(errE, w)
		return
	}

	errors.JSONSuccess(http.StatusOK, thread, w)
}
