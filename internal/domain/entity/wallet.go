package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	WalletTypeCommom WalletType = "commom"
	WalletTypeSeller WalletType = "seller"

	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeTransfer TransactionType = "transfer"

	EntryTypeInbound  = "inbound"
	EntryTypeOutbound = "outbound"
)

type (
	WalletType      string
	TransactionType string
	EntryType       string
	Transactions    []Transaction
)

func (t *Transactions) Balance() int64 {
	var amount int64
	for _, transaction := range *t {
		if transaction.EntryType == EntryTypeInbound {
			amount += int64(transaction.Amount)
			continue
		}
		amount -= int64(transaction.Amount)
	}
	return amount
}

type Wallet struct {
	ID              uuid.UUID    `json:"id"`
	CustomerName    string       `json:"customer_id"`
	WalletType      WalletType   `json:"wallet_type"`
	DocumentNumber  string       `json:"document_number"`
	Email           string       `json:"email"`
	EncodedPassword string       `json:"encoded_password"`
	Phone           string       `json:"phone"`
	Transactions    Transactions `json:"transactions"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func NewWallet(customerName, documentNumber, email, password, phone string) (*Wallet, error) {
	now := time.Now().UTC()
	wallet := Wallet{
		ID:             uuid.New(),
		CustomerName:   customerName,
		WalletType:     WalletTypeCommom,
		DocumentNumber: documentNumber,
		Email:          email,
		Phone:          phone,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return &wallet, nil
}

func (w *Wallet) Deposit(currency string, amount uint64) (*Transaction, error) {
	transactionID := uuid.New()
	t := Transaction{
		ID:              transactionID,
		ReferenceID:     transactionID,
		WalletID:        w.ID,
		TransactionType: TransactionTypeDeposit,
		EntryType:       EntryTypeInbound,
		Currency:        currency,
		Amount:          amount,
		Timestamp:       time.Now().UTC(),
	}

	return &t, t.Ok()
}

func (w *Wallet) Transfer(target Wallet, currency string, amount uint64) (*Transactions, error) {
	if w.WalletType == WalletTypeSeller {
		return nil, errors.New("TODO")
	}
	
	if w.Transactions.Balance() < int64(amount) {
		return nil, errors.New("TODO")
	}

	referenceID := uuid.New()
	now := time.Now().UTC()
	transactionSender := Transaction{
		ID:              uuid.New(),
		ReferenceID:     referenceID,
		Timestamp:       now,
		WalletID:        w.ID,
		TransactionType: TransactionTypeTransfer,
		EntryType:       EntryTypeOutbound,
		Currency:        currency,
		Amount:          amount,
	}

	transactionReceiver := Transaction{
		ID:              uuid.New(),
		ReferenceID:     referenceID,
		Timestamp:       now,
		WalletID:        target.ID,
		TransactionType: TransactionTypeTransfer,
		EntryType:       EntryTypeInbound,
		Currency:        currency,
		Amount:          amount,
	}

	return &Transactions{transactionReceiver, transactionSender}, nil
}

type Transaction struct {
	ID              uuid.UUID       `json:"id"`
	ReferenceID     uuid.UUID       `json:"reference_id"`
	WalletID        uuid.UUID       `json:"wallet_id"`
	TransactionType TransactionType `json:"transaction_type"`
	EntryType       EntryType       `json:"entry_type"`
	Currency        string          `json:"currency"`
	Amount          uint64          `json:"amount"`
	Snapshot        uuid.NullUUID   `json:"snapshot"`
	Timestamp       time.Time       `json:"timestamp"`
}

func (t *Transaction) Ok() error {
	return nil
}

func NewTransaction() (*Transaction, error) {
	t := Transaction{
		ID:        uuid.New(),
		Timestamp: time.Now().UTC(),
	}

	return &t, nil
}
