package usecase

import (
	"context"
	"time"

	"github.com/aburizalpurnama/go-simple-lending/internal/model"
	"github.com/aburizalpurnama/go-simple-lending/internal/payload/request"
	"github.com/aburizalpurnama/go-simple-lending/internal/repository"
	_payment "github.com/aburizalpurnama/go-simple-lending/pkg/payment"
	"gorm.io/gorm"
)

type (
	Payment interface {
		Create(ctx context.Context, accountId int, req request.CreatePayment) (model.Payment, error)
	}

	paymentImpl struct {
		db          *gorm.DB
		accountRepo repository.Account
		loanRepo    repository.Loan
		instRepo    repository.Installment
		paymentRepo repository.Payment
	}

	loanId        int
	installmentId int
)

func NewPayment(db *gorm.DB, accountRepo repository.Account, loanRepo repository.Loan, instRepo repository.Installment, paymentRepo repository.Payment) *paymentImpl {
	return &paymentImpl{
		db:          db,
		accountRepo: accountRepo,
		loanRepo:    loanRepo,
		instRepo:    instRepo,
		paymentRepo: paymentRepo,
	}
}

func (l *paymentImpl) Create(ctx context.Context, accountId int, req request.CreatePayment) (model.Payment, error) {
	payment := model.Payment{
		Amount:    req.Amount,
		Date:      time.Now().UTC(),
		AccountId: accountId,
	}

	err := l.db.Transaction(func(tx *gorm.DB) error {
		_, err := l.accountRepo.GetById(ctx, tx, accountId)
		if err != nil {
			return err
		}

		installments, err := l.instRepo.GetListAciveByAccountId(ctx, tx, accountId)
		if err != nil {
			return err
		}

		loanAdjustmentMap := map[int]int{}
		adjustedInstallmentSet := map[int]any{}

		adjustedLoanIds, err := _payment.AllocateAmount(ctx, installments, req.Amount, loanAdjustmentMap, adjustedInstallmentSet)
		if err != nil {
			return err
		}

		paymentId, err := l.paymentRepo.Create(ctx, tx, payment)
		if err != nil {
			return err
		}

		payment.Id = paymentId

		if len(adjustedLoanIds) > 0 {
			loans, err := l.loanRepo.GetListByIds(ctx, tx, adjustedLoanIds)
			if err != nil {
				return err
			}

			for _, loan := range loans {
				if amount, ok := loanAdjustmentMap[loan.Id]; ok {
					loan.PaidAmount += amount
					if loan.PaidAmount == loan.Amount {
						loan.Status = "paidoff"
					}

					err = l.loanRepo.UpdateNonZeroField(ctx, tx, loan)
					if err != nil {
						return err
					}
				}
			}
		}

		for _, inst := range installments {
			if _, ok := adjustedInstallmentSet[inst.Id]; ok {
				if inst.PaidAmount == inst.Amount {
					inst.Status = "paidoff"
				}

				err = l.instRepo.UpdateNonZeroField(ctx, tx, inst)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}
