package http

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("ThreadCreate")
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	//fmt.Println("ThreadCreate slug or id = ", slug)

	var posts []*models.Post
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		return
	}
	defer r.Body.Close()
	//fmt.Println("ThreadCreate posts = ", posts)

	resultPosts, errE := th.postUsecase.CreatePost(slug, posts)
	if errE != nil {
		//fmt.Println("ThreadCreate error = ", errE)
		errors.JSONError(errE, w)
		return
	}

	//fmt.Println("ThreadCreate result posts = ", resultPosts)
	errors.JSONSuccess(http.StatusCreated, resultPosts, w)
}

func (th *ThreadHandler) ThreadDetailsGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ThreadDetailsGET")
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	//fmt.Println("ThreadDetailsGET slug or id = ", slug)

	thread, errE := th.threadUsecase.GetThread(slug)
	if errE != nil {
		//fmt.Println("ThreadDetailsGET error = ", errE)
		errors.JSONError(errE, w)
		return
	}

	//fmt.Println("ThreadDetailsGET result thread = ", thread)
	errors.JSONSuccess(http.StatusOK, thread, w)
}

func (th *ThreadHandler) ThreadDetailsPOST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ThreadDetailsPOST")
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	//fmt.Println("ThreadDetailsPOST slug or id = ", slug)

	thread := &models.Thread{Forum: slug}
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		return
	}
	defer r.Body.Close()
	//fmt.Println("ThreadDetailsPOST thread = ", thread)

	resultThread, errE := th.threadUsecase.UpdateThread(thread)
	if errE != nil {
		//fmt.Println("ThreadDetailsPOST error = ", errE)
		errors.JSONError(errE, w)
		return
	}

	//fmt.Println("ThreadDetailsPOST result thread = ", resultThread)
	errors.JSONSuccess(http.StatusOK, resultThread, w)
}

func (th *ThreadHandler) ThreadPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ThreadPosts")
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	//fmt.Println("ThreadPosts slug or id = ", slug)

	postParams := &models.PostParams{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(postParams, r.URL.Query())
	if err != nil {
		return
	}
	defer r.Body.Close()
	//fmt.Println("ThreadPosts post params = ", postParams)

	posts, errE := th.postUsecase.GetPostsBySlugAndParams(slug, postParams)
	if errE != nil {
		//fmt.Println("ThreadPosts error = ", errE)
		errors.JSONError(errE, w)
		return
	}

	//fmt.Println("ThreadPosts result posts = ", postParams)
	errors.JSONSuccess(http.StatusOK, posts, w)
}

func (th *ThreadHandler) ThreadVote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ThreadVote")
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]
	//fmt.Println("ThreadVote slug or id = ", slug)

	vote := &models.Vote{}
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		return
	}
	defer r.Body.Close()
	//fmt.Println("ThreadVote vote = ", vote)

	thread, errE := th.threadUsecase.Vote(slug, vote)
	if errE != nil {
		//fmt.Println("ThreadVote error = ", errE)
		errors.JSONError(errE, w)
		return
	}

	//fmt.Println("ThreadVote result thread = ", thread)
	errors.JSONSuccess(http.StatusOK, thread, w)
}
