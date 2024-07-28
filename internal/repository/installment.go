package repository

import (
	"context"
	"net/http"

	"github.com/aburizalpurnama/go-simple-lending/internal/custerror"
	"github.com/aburizalpurnama/go-simple-lending/internal/model"
	_errors "github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Installment = new(installmentImpl)

type (
	Installment interface {
		Create(ctx context.Context, tx *gorm.DB, installment model.Installment) (int, error)
		CreateBulk(ctx context.Context, tx *gorm.DB, installments []model.Installment) error
		GetById(ctx context.Context, tx *gorm.DB, id int) (model.Installment, error)
		GetListByLoanId(ctx context.Context, tx *gorm.DB, loanId int) ([]model.Installment, error)
		GetListAciveByAccountId(ctx context.Context, tx *gorm.DB, accountId int) ([]model.Installment, error)
	}

	installmentImpl struct{}
)

func NewInstallment() *installmentImpl {
	return &installmentImpl{}
}

func (a *installmentImpl) Create(ctx context.Context, tx *gorm.DB, installment model.Installment) (int, error) {
	err := tx.Create(&installment).Error
	if err != nil {
		return 0, _errors.WithStack(err)
	}

	return installment.Id, nil
}

func (a *installmentImpl) CreateBulk(ctx context.Context, tx *gorm.DB, installments []model.Installment) error {
	if len(installments) > 0 {
		err := tx.CreateInBatches(installments, len(installments)).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *installmentImpl) GetById(ctx context.Context, tx *gorm.DB, id int) (model.Installment, error) {
	var installment model.Installment
	err := tx.Model(&model.Installment{}).Where("id", id).First(&installment).Error
	switch err {
	case nil:
		return installment, nil
	case gorm.ErrRecordNotFound:
		return model.Installment{}, custerror.New(http.StatusNotFound, "installment not found", err)
	default:
		return model.Installment{}, _errors.WithStack(err)
	}
}

func (a *installmentImpl) GetListByLoanId(ctx context.Context, tx *gorm.DB, loanId int) ([]model.Installment, error) {
	var installments []model.Installment
	err := tx.Model(&model.Installment{}).Where("loan_id", loanId).Find(&installments).Error
	if err != nil {
		return nil, _errors.WithStack(err)
	}

	return installments, nil
}

func (a *installmentImpl) GetListAciveByAccountId(ctx context.Context, tx *gorm.DB, accountId int) ([]model.Installment, error) {
	sqlQuery := `SELECT i.* 
	FROM installments i
	JOIN loans l ON i.loan_id = l.id
	WHERE 
		(i.amount - i.paid_amount) > 0
		AND l.account_id = ?
	ORDER By i.due_date ASC;`

	var installments []model.Installment
	err := tx.Raw(sqlQuery, accountId).Scan(&installments).Error
	if err != nil {
		return nil, _errors.WithStack(err)
	}

	return installments, nil
}
