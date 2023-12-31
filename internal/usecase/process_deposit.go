package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/internal/domain/entity"
	"github.com/guilhermealvess/guicpay/internal/domain/event"
	"github.com/guilhermealvess/guicpay/internal/domain/repository"
	"github.com/guilhermealvess/guicpay/internal/domain/service"
)

type processDeposit struct {
	repository        repository.WalletRepository
	eventNotification event.EventNotification
	authorizeService  service.AuthorizeService
}

func NewProcessDeposit(repo repository.WalletRepository, evnt event.EventNotification, auth service.AuthorizeService) ProcessDeposit {
	return &processDeposit{
		repository:        repo,
		eventNotification: evnt,
		authorizeService:  auth,
	}
}

type DepositParams struct {
	WalletReceiverID uuid.UUID
	Amount           uint64
	Currency         string
}

func (u processDeposit) Execute(ctx context.Context, params DepositParams) error {
	// ctx, cancel := context.WithTimeout(ctx, settings.Env.TransactionTimeout)
	// defer cancel()

	wallet, err := u.repository.GetWalletByID(ctx, params.WalletReceiverID)
	if err != nil {
		return err
	}

	if err := u.authorizeService.Authorize(ctx, wallet.Email, wallet.EncodedPassword); err != nil {
		return err
	}

	transaction, err := wallet.Deposit(params.Currency, params.Amount)
	if err != nil {
		return err
	}

	tx, err := u.repository.SaveTransactions(ctx, []entity.Transaction{*transaction})
	if err != nil {
		return err
	}

	if err := u.eventNotification.PublishEntity(ctx, []entity.Transaction{*transaction}); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
