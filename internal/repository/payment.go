package repository

import (
	"context"
	"net/http"

	"github.com/aburizalpurnama/go-simple-lending/internal/custerror"
	"github.com/aburizalpurnama/go-simple-lending/internal/model"
	_errors "github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Payment = new(paymentImpl)

type (
	Payment interface {
		Create(ctx context.Context, tx *gorm.DB, payment model.Payment) (int, error)
		GetById(ctx context.Context, tx *gorm.DB, id int) (model.Payment, error)
		GetListByAccountId(ctx context.Context, tx *gorm.DB, accountId int) ([]model.Payment, error)
	}

	paymentImpl struct{}
)

func NewPayment() *paymentImpl {
	return &paymentImpl{}
}

func (a *paymentImpl) Create(ctx context.Context, tx *gorm.DB, payment model.Payment) (int, error) {
	err := tx.Create(&payment).Error
	if err != nil {
		return 0, _errors.WithStack(err)
	}

	return payment.Id, nil
}

func (a *paymentImpl) GetById(ctx context.Context, tx *gorm.DB, id int) (model.Payment, error) {
	var payment model.Payment
	err := tx.Model(&model.Payment{}).Where("id", id).First(&payment).Error
	switch err {
	case nil:
		return payment, nil
	case gorm.ErrRecordNotFound:
		return model.Payment{}, custerror.New(http.StatusNotFound, "payment not found", err)
	default:
		return model.Payment{}, _errors.WithStack(err)
	}
}

func (a *paymentImpl) GetListByAccountId(ctx context.Context, tx *gorm.DB, accountId int) ([]model.Payment, error) {
	var payments []model.Payment
	err := tx.Model(&model.Payment{}).Where("account_id", accountId).Find(&payments).Error
	if err != nil {
		return nil, _errors.WithStack(err)
	}

	return payments, nil
}
