package models

import (
	"time"
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
