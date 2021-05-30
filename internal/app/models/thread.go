package models

import "time"

type Thread struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int       `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}

type ThreadWSlug struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int       `json:"votes"`
	Created time.Time `json:"created"`
}

func ConvertThread(thread *Thread) *ThreadWSlug {
	return &ThreadWSlug{
		thread.Id,
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Votes,
		thread.Created,
	}
}

type ThreadParams struct {
	Limit int    `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}

type Vote struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	ThreadId int    `json:"thread_id"`
}
