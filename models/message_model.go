package models

type Message struct {
	Id         string `json:"id"`
	CreateTime string `json:"create_time"`
	UserFrom   string `json:"user_from"`
	UserTo     string `json:"user_to"`
	Content    string `json:"content"`
	IsViewed   bool   `json:"is_viewed"`
}

type User struct {
	Id        int
	Email     string
	Password  string
	Icon      string
	Signature string
}
