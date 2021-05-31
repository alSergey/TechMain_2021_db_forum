package models

type Status struct {
	UserCount   int `json:"user"`
	ForumCount  int `json:"forum"`
	ThreadCount int `json:"thread"`
	PostCount   int `json:"post"`
}
