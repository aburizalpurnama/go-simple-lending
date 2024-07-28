package response

import "github.com/aburizalpurnama/go-simple-lending/internal/model"

type (
	Base struct {
		Description string `json:"description"`
		Data        any    `json:"data"`
	}

	GetAccount struct {
		model.Account
		AvailableLimit         int `json:"available_limit"`
		TotalLoanAmount        int `json:"total_loan_amount"`
		TotalPaidAmount        int `json:"total_paid_amount"`
		TotalOutstandingAmount int `json:"total_outstanding_amount"`
	}
)
