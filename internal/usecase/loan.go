package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/aburizalpurnama/go-simple-lending/internal/custerror"
	"github.com/aburizalpurnama/go-simple-lending/internal/model"
	"github.com/aburizalpurnama/go-simple-lending/internal/payload/request"
	"github.com/aburizalpurnama/go-simple-lending/internal/repository"
	"gorm.io/gorm"
)

type (
	Loan interface {
		Create(ctx context.Context, accountId int, req request.CreateLoan) (model.Loan, error)
	}

	loanImpl struct {
		db          *gorm.DB
		accountRepo repository.Account
		loanRepo    repository.Loan
		instRepo    repository.Installment
	}
)

func NewLoan(db *gorm.DB, accountRepo repository.Account, loanRepo repository.Loan, instRepo repository.Installment) *loanImpl {
	return &loanImpl{
		db:          db,
		accountRepo: accountRepo,
		loanRepo:    loanRepo,
		instRepo:    instRepo,
	}
}

func (l *loanImpl) Create(ctx context.Context, accountId int, req request.CreateLoan) (model.Loan, error) {
	loan := model.Loan{
		Amount:     req.Amount,
		PaidAmount: 0,
		Tenor:      req.Tenor,
		Date:       time.Now().UTC(),
		Status:     "active",
		AccountId:  accountId,
	}

	err := l.db.Transaction(func(tx *gorm.DB) error {
		account, err := l.accountRepo.GetById(ctx, tx, accountId)
		if err != nil {
			return err
		}

		osAmount, err := l.loanRepo.GetTotalOustandingByAccountId(ctx, tx, accountId)
		if err != nil {
			return err
		}

		if req.Amount > (account.Limit - osAmount) {
			return custerror.New(http.StatusBadRequest, "not enough limit", nil)
		}

		id, err := l.loanRepo.Create(ctx, tx, loan)
		if err != nil {
			return err
		}

		loan.Id = id

		installments := loan.GenerateInstallments()
		err = l.instRepo.CreateBulk(ctx, tx, installments)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return model.Loan{}, err
	}

	return loan, nil
}
