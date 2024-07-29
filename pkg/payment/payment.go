package payment

import (
	"context"
	"net/http"

	"github.com/aburizalpurnama/go-simple-lending/internal/custerror"
	"github.com/aburizalpurnama/go-simple-lending/internal/model"
)

// AllocateAmount allocating payment amount to active installments and loans.
// loanAdjsutmentMap and adjustedInstallmentSet shall a non-nil value.
// It returns adjusted loan ids and error when payment amount bigger than outstanding.
func AllocateAmount(ctx context.Context, installments []model.Installment, paymentAmount int, loanAdjustmentMap map[int]int, adjustedInstallmentSet map[int]any) (adjustedLoanIds []int, err error) {
	if loanAdjustmentMap == nil {
		loanAdjustmentMap = map[int]int{}
	}

	if adjustedInstallmentSet == nil {
		adjustedInstallmentSet = map[int]any{}
	}

	remainingAmount := paymentAmount
	for i := range installments {
		osAmount := (installments[i].Amount - installments[i].PaidAmount)
		if remainingAmount > 0 && osAmount > 0 {
			adjustmentAmount := osAmount
			if remainingAmount < osAmount {
				adjustmentAmount = remainingAmount
			}

			if _, ok := loanAdjustmentMap[installments[i].LoanId]; ok {
				loanAdjustmentMap[installments[i].LoanId] += adjustmentAmount
			} else {
				loanAdjustmentMap[installments[i].LoanId] = adjustmentAmount
				adjustedLoanIds = append(adjustedLoanIds, installments[i].LoanId)
			}

			installments[i].PaidAmount += adjustmentAmount

			adjustedInstallmentSet[installments[i].Id] = nil
			remainingAmount -= adjustmentAmount
		}
	}

	if remainingAmount > 0 {
		return nil, custerror.New(http.StatusBadRequest, "too much payment amount", nil)
	}

	return
}
