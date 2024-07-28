package model

import "time"

type Loan struct {
	Id         int       `json:"id"`
	Amount     int       `json:"amount"`
	PaidAmount int       `json:"paid_amount"`
	Date       time.Time `json:"date"`
	Status     string    `json:"status"`

	AccountId int `json:"account_id"`
}
