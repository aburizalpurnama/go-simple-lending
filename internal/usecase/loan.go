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
		laonRepo    repository.Loan
	}
)

func NewLoan(db *gorm.DB, accountRepo repository.Account, loanRepo repository.Loan) *loanImpl {
	return &loanImpl{
		db:          db,
		accountRepo: accountRepo,
		laonRepo:    loanRepo,
	}
}

func (l *loanImpl) Create(ctx context.Context, accountId int, req request.CreateLoan) (model.Loan, error) {
	loan := model.Loan{
		Amount:     req.Amount,
		PaidAmount: 0,
		Date:       time.Now().UTC(),
		Status:     "active",
		AccountId:  accountId,
	}

	err := l.db.Transaction(func(tx *gorm.DB) error {
		account, err := l.accountRepo.GetById(ctx, tx, accountId)
		if err != nil {
			return err
		}

		osAmount, err := l.laonRepo.GetTotalOustandingByAccountId(ctx, tx, accountId)
		if err != nil {
			return err
		}

		if req.Amount > (account.Limit - osAmount) {
			return custerror.New(http.StatusBadRequest, "not enough limit", nil)
		}

		id, err := l.laonRepo.Create(ctx, tx, loan)
		if err != nil {
			return err
		}

		loan.Id = id

		return nil
	})
	if err != nil {
		return model.Loan{}, err
	}

	return loan, nil
}
