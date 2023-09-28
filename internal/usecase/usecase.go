package usecase

import "context"

type RegisterWallet interface {
	Execute(ctx context.Context, params RegisterWalletParams) error
}

type ProcessDeposit interface {
	Execute(ctx context.Context, params DepositParams) error
}

type ProcessTransfer interface {
	Execute(ctx context.Context, params ProcessTransferParams) error
}
