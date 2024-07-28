package repository

import (
	"context"
	"net/http"

	"github.com/aburizalpurnama/go-simple-lending/internal/custerror"
	"github.com/aburizalpurnama/go-simple-lending/internal/model"
	_errors "github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Account = new(accountImpl)

type (
	Account interface {
		Create(ctx context.Context, tx *gorm.DB, account model.Account) (int, error)
		GetById(ctx context.Context, tx *gorm.DB, id int) (model.Account, error)
	}

	accountImpl struct{}
)

func NewAccount() *accountImpl {
	return &accountImpl{}
}

func (a *accountImpl) Create(ctx context.Context, tx *gorm.DB, account model.Account) (int, error) {
	err := tx.Create(&account).Error
	if err != nil {
		return 0, _errors.WithStack(err)
	}

	return account.Id, nil
}

func (a *accountImpl) GetById(ctx context.Context, tx *gorm.DB, id int) (model.Account, error) {
	var account model.Account
	err := tx.Model(&model.Account{}).Where("id", id).First(&account).Error
	switch err {
	case nil:
		return account, nil
	case gorm.ErrRecordNotFound:
		return model.Account{}, custerror.New(http.StatusNotFound, "account not found", err)
	default:
		return model.Account{}, _errors.WithStack(err)
	}
}
