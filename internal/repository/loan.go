package repository

import (
	"context"
	"net/http"

	"github.com/aburizalpurnama/go-simple-lending/internal/custerror"
	"github.com/aburizalpurnama/go-simple-lending/internal/model"
	_errors "github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Loan = new(loanImpl)

type (
	Loan interface {
		Create(ctx context.Context, tx *gorm.DB, loan model.Loan) (int, error)
		GetById(ctx context.Context, tx *gorm.DB, id int) (model.Loan, error)
		GetListByAccountId(ctx context.Context, tx *gorm.DB, accountId int) ([]model.Loan, error)
		GetTotalOustandingByAccountId(ctx context.Context, tx *gorm.DB, accountId int) (osAmount int, err error)
	}

	loanImpl struct{}
)

func NewLoan() *loanImpl {
	return &loanImpl{}
}

func (a *loanImpl) Create(ctx context.Context, tx *gorm.DB, loan model.Loan) (int, error) {
	err := tx.Create(&loan).Error
	if err != nil {
		return 0, _errors.WithStack(err)
	}

	return loan.Id, nil
}

func (a *loanImpl) GetById(ctx context.Context, tx *gorm.DB, id int) (model.Loan, error) {
	var loan model.Loan
	err := tx.Model(&model.Loan{}).Where("id", id).First(&loan).Error
	switch err {
	case nil:
		return loan, nil
	case gorm.ErrRecordNotFound:
		return model.Loan{}, custerror.New(http.StatusNotFound, "loan not found", err)
	default:
		return model.Loan{}, _errors.WithStack(err)
	}
}

func (a *loanImpl) GetListByAccountId(ctx context.Context, tx *gorm.DB, accountId int) ([]model.Loan, error) {
	var loans []model.Loan
	err := tx.Model(&model.Loan{}).Where("account_id", accountId).Find(&loans).Error
	if err != nil {
		return nil, _errors.WithStack(err)
	}

	return loans, nil
}

func (a *loanImpl) GetTotalOustandingByAccountId(ctx context.Context, tx *gorm.DB, accountId int) (osAmount int, err error) {
	err = tx.Model(&model.Loan{}).Select("COALESCE(SUM(amount) - SUM(paid_amount), 0)").Where("account_id", accountId).Find(&osAmount).Error
	if err != nil {
		return 0, _errors.WithStack(err)
	}

	return osAmount, nil
}
