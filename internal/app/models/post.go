package models

import "time"

type Post struct {
	Id       int       `json:"id"`
	Parent   int64     `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"is_edited"`
	Forum    string    `json:"forum"`
	Thread   int       `json:"thread"`
	Created  time.Time `json:"created"`
}
