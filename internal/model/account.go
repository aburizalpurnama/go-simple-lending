package model

type Account struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Limit int    `json:"limit"`
}
