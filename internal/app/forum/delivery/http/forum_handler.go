package http

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/forum"
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

}

func (fh *ForumHandler) ForumDetails(w http.ResponseWriter, r *http.Request) {

}

func (fh *ForumHandler) ForumCreateThread(w http.ResponseWriter, r *http.Request) {

}

func (fh *ForumHandler) ForumUsers(w http.ResponseWriter, r *http.Request) {

}

func (fh *ForumHandler) ForumThreads(w http.ResponseWriter, r *http.Request) {

}
