package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/internal/domain/event"
	"github.com/guilhermealvess/guicpay/internal/domain/repository"
	"github.com/guilhermealvess/guicpay/internal/domain/service"
)

type processTransfer struct {
	repository        repository.WalletRepository
	eventNotification event.EventNotification
	authorizeService  service.AuthorizeService
	mutexService      service.MutexService
}

type ProcessTransferParams struct {
	SenderWalletID   uuid.UUID `json:"sender_wallet_id"`
	ReceiverWalletID uuid.UUID `json:"receiver_wallet_id"`
	Amount           uint64    `json:"amount"`
	Currency         string    `json:"currency"`
}

func (u processTransfer) Execute(ctx context.Context, params ProcessTransferParams) error {
	sender, err := u.repository.GetWalletByID(ctx, params.SenderWalletID)
	if err != nil {
		return err
	}

	receiver, err := u.repository.GetWalletByID(ctx, params.ReceiverWalletID)
	if err != nil {
		return err
	}

	transactions, err := sender.Transfer(*receiver, params.Currency, params.Amount)
	if err != nil {
		return err
	}

	mutex := u.mutexService.NewMutex(sender.ID.String(), 7*time.Second)
	if err := mutex.Lock(ctx); err != nil {
		return err
	}
	defer mutex.Unlock(ctx)

	tx, err := u.repository.SaveTransactions(ctx, *transactions)
	if err != nil {
		return err
	}

	if err := u.authorizeService.Authorize(ctx, sender.Email, sender.EncodedPassword); err != nil {
		tx.Rollback()
		return err
	}

	if err := u.eventNotification.PublishTransactions(ctx, *transactions); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
