package model

import "time"

type Installment struct {
	Id         int       `json:"id"`
	Amount     int       `json:"amount"`
	PaidAmount int       `json:"paid_amount"`
	DueDate    time.Time `json:"due_date"`
	Status     string    `json:"status"`

	LoanId int `json:"loan_id"`
}
