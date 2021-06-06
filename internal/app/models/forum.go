package models

//easyjson:json
type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}

//easyjson:json
type ForumParams struct {
	Limit int    `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}