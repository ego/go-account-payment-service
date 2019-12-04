// dal.go data access layer (DAL) for application.
package account

import (
	"context"
	"database/sql"
	"github.com/go-kit/kit/log"
)

// dal - application DAL definition with database connection and logger.
type dal struct {
	db     *sql.DB
	logger log.Logger
}

// NewDal - provide ability to create Dal.
func NewDal(db *sql.DB, logger log.Logger) Dal {
	return &dal{
		db:     db,
		logger: log.With(logger, "dal", "sql"),
	}
}

// CreateAccount Dal method - interact with database and create new row in DB table account.
func (dal *dal) CreateAccount(ctx context.Context, account Account) (Account, error) {
	sql := `
		INSERT INTO account (name, balance, currency)
		VALUES ($1, $2, $3) RETURNING account.*`

	// Validate account balance, it must be > 0
	if account.Balance.LessThan(ZERO) {
		return Account{}, AccountAmountErr
	}

	// Try create account in DB
	newAccount := Account{}
	err := dal.db.QueryRow(sql, account.Name, account.Balance, account.Currency).Scan(
		&newAccount.ID,
		&newAccount.Name,
		&newAccount.Balance,
		&newAccount.Currency,
		&newAccount.CreatedAt,
	)
	if err != nil {
		return Account{}, err
	}

	return newAccount, nil
}

// GetAccount Dal method - interact with database and retrieve account row from DB table.
func (dal *dal) GetAccount(ctx context.Context, id int) (Account, error) {
	account := Account{}
	// Select query for account
	err := dal.db.QueryRow("SELECT id, name, balance, currency, created_at FROM account WHERE id=$1", id).Scan(
		&account.ID,
		&account.Name,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)
	// Handle some errors and return custom more API friendly error
	if err != nil {
		dal.logger.Log("Error", err)
		switch err {
		case sql.ErrNoRows:
			return Account{}, NoResultErr
		default:
			return Account{}, err
		}
		return Account{}, err
	}
	return account, nil
}

// selectForUpdate - helper function that provide SELECT FOR UPDATE postgresql feature.
func selectForUpdate(tx *sql.Tx, id int) (Account, error) {
	account := Account{}

	sql := `SELECT id, name, balance, currency, created_at FROM account WHERE id=$1 FOR UPDATE`

	err := tx.QueryRow(sql, id).Scan(
		&account.ID,
		&account.Name,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)
	return account, err
}

// makeAccountPayment - helper function create AccountPayment struct
// and determine payment direction.
func makeAccountPayment(payment Payment, account1 Account, account2 Account) AccountPayment {
	ap := AccountPayment{TxPayment: payment}
	if ap.TxPayment.Direction == OUTGOING {
		ap.CreditAccount = account1
		ap.DepositAccount = account2
	} else {
		ap.CreditAccount = account2
		ap.DepositAccount = account1
	}
	return ap
}

// updateAccounts - helper function provide UPDATE postgresql feature.
// Updates two account rows in DB.
func updateAccounts(tx *sql.Tx, ap AccountPayment) (error, error) {
	update1 := `UPDATE account SET balance = balance - $2 WHERE id=$1`
	update2 := `UPDATE account SET balance = balance + $2 WHERE id=$1`
	_, err1 := tx.Exec(update1, ap.CreditAccount.ID, ap.TxPayment.Amount)
	_, err2 := tx.Exec(update2, ap.DepositAccount.ID, ap.TxPayment.Amount)
	return err1, err2
}

// insertPayment - helper function provide INSERT postgresql feature.
// Inserts payment into DB.
func insertPayment(tx *sql.Tx, ap AccountPayment) (Payment, error) {
	insert := `
		INSERT INTO payment (account_id, to_account_id, amount, direction)
		VALUES ($1, $2, $3, $4) RETURNING payment.*`

	newPayment := Payment{}
	err := tx.QueryRow(insert, ap.TxPayment.AccountID, ap.TxPayment.ToAccountID, ap.TxPayment.Amount, ap.TxPayment.Direction).Scan(
		&newPayment.ID,
		&newPayment.AccountID,
		&newPayment.ToAccountID,
		&newPayment.Amount,
		&newPayment.Direction,
		&newPayment.CreatedAt,
	)
	return newPayment, err
}

// CreatePayment Dal method - interact with database and
// create new row in DB table payment.
func (dal *dal) CreatePayment(ctx context.Context, payment Payment) (Payment, error) {
	// Validate source account and destination account
	if payment.AccountID == payment.ToAccountID {
		return Payment{}, SameAccountErr
	}

	// Begin transaction
	tx, err := dal.db.Begin()
	if err != nil {
		return Payment{}, err
	}

	// Write and read locks first account
	account1, err := selectForUpdate(tx, payment.AccountID)
	if err != nil {
		// If we can not get row lock for first account or timeout happens
		// rollback database transaction
		tx.Rollback()
		dal.logger.Log(err)
		return Payment{}, err
	}

	// Write and read locks second account
	account2, err := selectForUpdate(tx, payment.ToAccountID)
	if err != nil {
		// If we can not get row lock for second account or timeout happens
		// rollback database transaction
		tx.Rollback()
		dal.logger.Log(err)
		return Payment{}, err
	}

	// Determine payment direction
	ap := makeAccountPayment(payment, account1, account2)

	// Validate same currency
	if ap.CreditAccount.Currency != ap.DepositAccount.Currency {
		return Payment{}, AccountCurrencyErr
	}

	// Update credit and debit accounts
	err1, err2 := updateAccounts(tx, ap)
	if err1 != nil || err2 != nil {
		// If any error on update account table rollback database transaction
		tx.Rollback()
		dal.logger.Log(err1)
		dal.logger.Log(err2)
		if err1 != nil {
			return Payment{}, err1
		} else {
			return Payment{}, err2
		}
	}

	// Insert new payment
	newPayment, err := insertPayment(tx, ap)
	if err != nil {
		// If something goes wrong on insert new payment into database
		// rollback database transaction
		tx.Rollback()
		dal.logger.Log(err)
		return Payment{}, err
	}

	// End transaction
	err = tx.Commit()
	if err != nil {
		// If something goes wrong on commit rollback database transaction
		tx.Rollback()
		dal.logger.Log(err)
		return Payment{}, TransactionErr
	}

	return newPayment, nil
}

// GetPayment Dal method - interact with database and retrieve payment row from DB table.
func (dal *dal) GetPayment(ctx context.Context, id int) (Payment, error) {
	payment := Payment{}
	err := dal.db.QueryRow("SELECT id, account_id, to_account_id, amount, direction, created_at FROM payment WHERE id=$1", id).Scan(
		&payment.ID,
		&payment.AccountID,
		&payment.ToAccountID,
		&payment.Amount,
		&payment.Direction,
		&payment.CreatedAt,
	)

	// handle some errors and return custom more API friendly error.
	if err != nil {
		dal.logger.Log("Error", err)
		switch err {
		case sql.ErrNoRows:
			return Payment{}, NoResultErr
		default:
			return Payment{}, err
		}
		return Payment{}, err
	}
	return payment, nil
}
