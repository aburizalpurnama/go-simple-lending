package model

import "time"

type Loan struct {
	Id         int       `json:"id"`
	Amount     int       `json:"amount"`
	PaidAmount int       `json:"paid_amount"`
	Tenor      int       `json:"tenor"`
	Date       time.Time `json:"date"`
	Status     string    `json:"status"`

	AccountId int `json:"account_id"`
}

func (l *Loan) GenerateInstallments() []Installment {
	installments := make([]Installment, l.Tenor)

	instAmounts := l.getInstallmentAmounts()
	for i := range installments {
		installments[i].Amount = instAmounts[i]
		installments[i].PaidAmount = 0
		installments[i].DueDate = l.Date.AddDate(0, i+1, 0)
		installments[i].Status = "active"
		installments[i].LoanId = l.Id
	}

	return installments
}

func (l *Loan) getInstallmentAmounts() []int {
	mod := l.Amount % l.Tenor

	instAmounts := make([]int, l.Tenor)
	for i := range instAmounts {
		instAmounts[i] = l.Amount / l.Tenor
		if mod > 0 {
			instAmounts[i]++
			mod--
		}

	}

	return instAmounts
}
