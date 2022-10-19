package models

type SearchRes struct {
	ID    int32  `json:"id" db:"id"`
	URL   string `json:"url" db:"url"`
	Title string `json:"title" db:"title"`
}

type SearchQuery struct {
	Text string `json:"text"`
}
