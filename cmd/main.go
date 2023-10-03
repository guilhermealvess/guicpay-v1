package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Rhymond/go-money"
	"github.com/guilhermealvess/guicpay/internal/infra/event"
	"github.com/guilhermealvess/guicpay/internal/infra/repository"
	"github.com/guilhermealvess/guicpay/internal/infra/service"
	"github.com/guilhermealvess/guicpay/internal/settings"
	"github.com/guilhermealvess/guicpay/internal/usecase"
	_ "github.com/lib/pq"
)

func fakeWallet() usecase.RegisterWalletParams {
	const data = `{
		"customer_name": "Guilherme Alves",
		"document_number": "0123456789-10",
		"email": "guilhermeasilva.dev@gmail.com",
		"password": "guicpay-simplificado",
		"phone": "+5534996344108",
		"wallet_type": "commom"
	}`
	var params usecase.RegisterWalletParams
	if err := json.Unmarshal([]byte(data), &params); err != nil {
		panic(err)
	}

	return params
}

func main() {
	db := NewDatabaseConnection()
	eventNotification := event.NewEventNotification(settings.Env.BrokerStreamUrl)
	authorizeService := service.NewAuthorizeService(settings.Env.AuthorizeServiceUrl)
	mutexService := service.NewMutexService(settings.Env.RedisUrl)
	walletRepository := repository.NewWalletRepository(db)

	registerWalletUsecase := usecase.NewRegisterWallet(walletRepository, eventNotification, authorizeService)
	processDepositUsecase := usecase.NewProcessDeposit(walletRepository, eventNotification, authorizeService)
	processTransferUsecase := usecase.NewProcessTransfer(walletRepository, eventNotification, authorizeService, mutexService)
	calcBalanceUsecase := usecase.NewCalculatorBalance(walletRepository)

	ctx := context.Background()
	params := fakeWallet()

	wallet, err := registerWalletUsecase.Execute(ctx, usecase.RegisterWalletParams{
		CustomerName:   params.CustomerName,
		DocumentNumber: params.DocumentNumber,
		Email:          params.Email,
		Password:       params.Password,
		Phone:          params.Phone,
		WalletType:     params.WalletType,
	})
	if err != nil {
		panic(err)
	}

	walletAux, err := registerWalletUsecase.Execute(ctx, usecase.RegisterWalletParams{
		CustomerName:   "Fulano De Tal",
		DocumentNumber: "123456789-11",
		Email:          "fulano@gmail.com",
		Password:       params.Password,
		Phone:          "+5534988087707",
		WalletType:     params.WalletType,
	})
	if err != nil {
		panic(err)
	}

	if err := processDepositUsecase.Execute(ctx, usecase.DepositParams{
		WalletReceiverID: wallet.ID,
		Amount:           100000 * 100,
		Currency:         money.BRL,
	}); err != nil {
		panic(err)
	}

	if err := processTransferUsecase.Execute(ctx, usecase.ProcessTransferParams{
		SenderWalletID:   wallet.ID,
		ReceiverWalletID: walletAux.ID,
		Amount:           7521,
		Currency:         money.BRL,
	}); err != nil {
		panic(err)
	}

	balance, err := calcBalanceUsecase.Execute(ctx, wallet.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Balance: %d\n", balance)
}

func NewDatabaseConnection() *sql.DB {
	db, err := sql.Open("postgres", settings.Env.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
