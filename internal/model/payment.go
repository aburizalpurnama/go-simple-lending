package model

import "time"

type Payment struct {
	Id     int       `json:"id"`
	Amount int       `json:"amount"`
	Date   time.Time `json:"date"`

	AccountId int `json:"account_id"`
}
