package usecase

import (
	"context"

	"github.com/guilhermealvess/guicpay/internal/domain/entity"
	"github.com/guilhermealvess/guicpay/internal/domain/event"
	"github.com/guilhermealvess/guicpay/internal/domain/repository"
	"github.com/guilhermealvess/guicpay/internal/domain/service"
)

type RegisterWalletParams struct {
	CustomerName   string            `json:"customer_name"`
	DocumentNumber string            `json:"document_number"`
	Email          string            `json:"email"`
	Password       string            `json:"password"`
	Phone          string            `json:"phone"`
	WalletType     entity.WalletType `json:"wallet_type"`
}

type registerWallet struct {
	repository        repository.WalletRepository
	eventNotification event.EventNotification
	authorizeService  service.AuthorizeService
}

func NewRegisterWallet(
	repo repository.WalletRepository,
	evntNotify event.EventNotification,
	auth service.AuthorizeService,
) RegisterWallet {
	return &registerWallet{
		repository:        repo,
		eventNotification: evntNotify,
		authorizeService:  auth,
	}
}

func (u registerWallet) Execute(ctx context.Context, params RegisterWalletParams) (*entity.Wallet, error) {
	wallet, err := entity.NewWallet(params.CustomerName, params.DocumentNumber, params.Email, params.Password, params.Phone)
	if err != nil {
		return nil, err
	}

	tx, err := u.repository.SaveWallet(ctx, *wallet)
	if err != nil {
		return nil, err
	}

	if err := u.authorizeService.RegisterUser(ctx, wallet.ID.String(), wallet.Email, wallet.EncodedPassword); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := u.eventNotification.PublishEntity(ctx, *wallet); err != nil {
		tx.Rollback()
		return nil, err
	}

	return wallet, tx.Commit()
}
