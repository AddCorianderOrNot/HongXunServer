package models

type Message struct {
	Id         int
	CreateTime int
	UserFrom   int
	UserTo     int
	Content    string
	IsViewed   bool
}
