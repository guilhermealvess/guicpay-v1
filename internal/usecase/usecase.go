package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/internal/domain/entity"
)

type RegisterWallet interface {
	Execute(ctx context.Context, params RegisterWalletParams) (*entity.Wallet, error)
}

type ProcessDeposit interface {
	Execute(ctx context.Context, params DepositParams) error
}

type ProcessTransfer interface {
	Execute(ctx context.Context, params ProcessTransferParams) error
}

type CalculatorBalance interface {
	Execute(ctx context.Context, walletID uuid.UUID) (int64, error)
}
