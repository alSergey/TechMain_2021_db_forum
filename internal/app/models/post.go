package models

import (
	"time"

	"github.com/jackc/pgx/pgtype"
)

type Post struct {
	Id       int              `json:"id"`
	Parent   int64            `json:"parent"`
	Author   string           `json:"author"`
	Message  string           `json:"message"`
	IsEdited bool             `json:"is_edited"`
	Forum    string           `json:"forum"`
	Thread   int              `json:"thread"`
	Created  time.Time        `json:"created"`
	Path     pgtype.Int8Array `json:"-"`
}

type PostParams struct {
	Limit int    `json:"limit"`
	Since int    `json:"since"`
	Sort  string `json:"sort"`
	Desc  bool   `json:"desc"`
}
