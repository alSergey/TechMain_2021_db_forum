package models

import (
	"time"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/uuid"
)

type Post struct {
	Id       int       `json:"id"`
	Parent   int64     `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   int       `json:"thread"`
	Created  time.Time `json:"created"`
}

type PostParams struct {
	Limit int    `json:"limit"`
	Since int    `json:"since"`
	Sort  string `json:"sort"`
	Desc  bool   `json:"desc"`
}

type FullPostParams struct {
	User   bool `json:"user"`
	Forum  bool `json:"forum"`
	Thread bool `json:"thread"`
}

type FullPost struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author"`
	Forum  *Forum  `json:"forum"`
	Thread *Thread `json:"thread"`
}

type FullPostWSlug struct {
	Post   *Post        `json:"post"`
	Author *User        `json:"author"`
	Forum  *Forum       `json:"forum"`
	Thread *ThreadWSlug `json:"thread"`
}

func GetResultPost(post *FullPost) interface{} {
	var result interface{}
	result = post

	if post.Thread != nil {
		if uuid.IsCreatedSlug(post.Thread.Slug) {
			result = &FullPostWSlug{
				Post:   post.Post,
				Author: post.Author,
				Forum:  post.Forum,
				Thread: ConvertThread(post.Thread),
			}
		}
	} else {
		result = post
	}

	return result
}

type GetPostType uint8

const (
	GetPost GetPostType = iota
	GetUser
	GetThread
	GetForum
	GetUserThread
	GetUserForum
	GetThreadForum
	GetUserThreadForum
)
