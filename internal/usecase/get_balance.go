package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/internal/domain/repository"
)

type calculatorBalance struct {
	walletRepository repository.WalletRepository
}

func NewCalculatorBalance(repo repository.WalletRepository) CalculatorBalance {
	return &calculatorBalance{repo}
}

func (u calculatorBalance) Execute(ctx context.Context, walletID uuid.UUID) (int64, error) {
	wallet, err := u.walletRepository.GetWalletByID(ctx, walletID)
	if err != nil {
		return 0, err
	}

	return wallet.Transactions.Balance(), nil
}
