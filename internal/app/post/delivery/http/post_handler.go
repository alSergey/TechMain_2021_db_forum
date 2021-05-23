package http

import (
	"net/http"

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

}

func (ph *PostHandler) PostDetailsPOST(w http.ResponseWriter, r *http.Request) {

}
