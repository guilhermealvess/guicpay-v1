package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/internal/domain/entity"
	"github.com/guilhermealvess/guicpay/internal/domain/repository"
	"github.com/guilhermealvess/guicpay/internal/sql/queries"
)

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository() *walletRepository {
	return &walletRepository{}
}

func (r walletRepository) SaveTransaction(ctx context.Context, transactions []entity.Transaction) (repository.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := queries.New(tx)

	for _, t := range transactions {
		params := queries.SaveTransactionParams{
			ID:          t.ID,
			ReferenceID: t.ReferenceID,
			WalletID:    t.WalletID,
			Type:        queries.TransactionType(t.TransactionType),
			EntriesType: queries.TransactionEntryType(t.EntryType),
			Currency:    t.Currency,
			Amount:      int32(t.Amount),
			Timestamp:   t.Timestamp,
		}
		if err := query.SaveTransaction(ctx, params); err != nil {
			return nil, err
		}
	}

	return tx, nil
}

func (r walletRepository) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*entity.Wallet, error) {
	query := queries.New(r.db)
	rows, err := query.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	var wallet entity.Wallet
	transactions := make([]entity.Transaction, len(rows))
	for i, row := range rows {
		if i == 0 {
			wallet.ID = row.WalletID
			wallet.CustomerName = row.CustomerName
			wallet.DocumentNumber = row.DocumentNumber
			wallet.Email = row.Email
			wallet.EncodedPassword = row.EncodedPassword
			wallet.Phone = row.PhoneNumber
			wallet.CreatedAt = row.CreatedAt
			wallet.UpdatedAt = row.UpdatedAt
			wallet.WalletType = entity.WalletType(row.WalletType)
		}

		t := entity.Transaction{
			ID:              row.TransactionID,
			ReferenceID:     row.ReferenceID,
			WalletID:        row.WalletID,
			TransactionType: entity.TransactionType(row.TransactionType),
			EntryType:       entity.EntryType(row.EntriesType),
			Currency:        row.Currency,
			Amount:          uint64(row.Amount),
			Timestamp:       row.Timestamp,
		}
		transactions[i] = t
	}
	wallet.Transactions = transactions

	return &wallet, nil
}

func (r walletRepository) SaveTransactions(ctx context.Context, transactions []entity.Transaction) (repository.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := queries.New(tx)

	for _, t := range transactions {
		params := queries.SaveTransactionParams{
			ID:          t.ID,
			ReferenceID: t.ReferenceID,
			WalletID:    t.WalletID,
			Type:        queries.TransactionType(t.TransactionType),
			EntriesType: queries.TransactionEntryType(t.EntryType),
			Timestamp:   t.Timestamp,
			Currency:    t.Currency,
			Amount:      int32(t.Amount),
		}
		if err := query.SaveTransaction(ctx, params); err != nil {
			return nil, err
		}
	}

	return tx, nil
}
