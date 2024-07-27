package request

type (
	CreateAccount struct {
		Name  string `validate:"require" json:"name"`
		Limit int    `validate:"require" json:"limit"`
	}

	CreateLoan struct {
		Limit int `validate:"require" json:"amount"`
	}

	CreatePayment struct {
		Amount int `validate:"require" json:"amount"`
	}
)
