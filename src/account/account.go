// Package implements WEB service API for account and payment features.
package account

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
)

// Account representations, keeps account data.
type Account struct {
	ID        int             `json:"id,omitempty"`
	Name      string          `json:"name"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	CreatedAt string          `json:"created_at,omitempty"`
}

// Payment representations, keeps Payment data.
type Payment struct {
	ID          int             `json:"id,omitempty"`
	AccountID   int             `json:"account_id"`
	ToAccountID int             `json:"to_account_id"`
	Amount      decimal.Decimal `json:"amount"`
	Direction   string          `json:"direction"`
	CreatedAt   string          `json:"created_at,omitempty"`
}

// AccountPayment representations, keeps two Accounts and Payment data.
type AccountPayment struct {
	CreditAccount  Account
	DepositAccount Account
	TxPayment      Payment
}

// Enum
const OUTGOING = "outgoing"

// Base vars and errors.
// BalanceErr - errors related to account balance.
// NoResultErr - errors related to database and empty results.
// SameAccountErr - errors related to account and business logic.
// AccountAmountErr - errors related to account amount and business logic.
// AccountCurrencyErr - errors related to account currency and business logic.
// TransactionErr - errors related to database and concurrent updates.
var (
	ZERO, _            = decimal.NewFromString("0.0")
	BalanceErr         = errors.New("Error, not enough balance!")
	NoResultErr        = errors.New("Error, no rows in result set!")
	SameAccountErr     = errors.New("Error, the account value must be different from the to_account value!")
	AccountAmountErr   = errors.New("Error, the amount must be greater than 0!")
	AccountCurrencyErr = errors.New("Error, accounts must be the same currency!")
	TransactionErr     = errors.New("Error, account payment transaction conflict!")
)

// Dal interface - provide data and database opperations.
type Dal interface {
	CreateAccount(ctx context.Context, user Account) (Account, error)
	GetAccount(ctx context.Context, id int) (Account, error)

	CreatePayment(ctx context.Context, payment Payment) (Payment, error)
	GetPayment(ctx context.Context, id int) (Payment, error)
}
