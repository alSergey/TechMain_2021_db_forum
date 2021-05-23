package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/thread"
)

type ThreadHandler struct {
	threadUsecase thread.ThreadUsecase
}

func NewThreadHandler(threadUsecase thread.ThreadUsecase) *ThreadHandler {
	return &ThreadHandler{
		threadUsecase: threadUsecase,
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

}

func (th *ThreadHandler) ThreadDetailsGET(w http.ResponseWriter, r *http.Request) {

}

func (th *ThreadHandler) ThreadDetailsPOST(w http.ResponseWriter, r *http.Request) {

}

func (th *ThreadHandler) ThreadPosts(w http.ResponseWriter, r *http.Request) {

}

func (th *ThreadHandler) ThreadVote(w http.ResponseWriter, r *http.Request) {

}
