package models

type Option struct {
	Text string `json:"text"`
}

type Survey struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Options []Option `json:"options"`
}