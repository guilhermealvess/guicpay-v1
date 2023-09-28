package event

import (
	"context"

	"github.com/guilhermealvess/guicpay/internal/domain/entity"
)

type EventNotification interface {
	PublishWallet(ctx context.Context, wallet entity.Wallet) error
	PublishTransactions(ctx context.Context, transactions entity.Transactions) error
}