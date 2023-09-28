package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/internal/domain/entity"
)

type WalletRepository interface {
	SaveWallet(ctx context.Context, wallet entity.Wallet) (Tx, error)
	GetWalletByID(ctx context.Context, id uuid.UUID) (*entity.Wallet, error)

	SaveTransactions(ctx context.Context, transactions []entity.Transaction) (Tx, error)
}

type Tx interface {
	Commit() error
	Rollback() error
}
