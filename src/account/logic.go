// logic.go provide base logic for application.
package account

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// service - application service definition with data access layer (DAL) and logger.
type service struct {
	db     Dal
	logger log.Logger
}

// NewService - provide ability to create Service.
func NewService(pg Dal, logger log.Logger) Service {
	return &service{
		db:     pg,
		logger: logger,
	}
}

// CreateAccount API method for Service provide create account feature.
func (s service) CreateAccount(ctx context.Context, account Account) (Account, error) {
	logger := log.With(s.logger, "method", "CreateAccount")

	account, err := s.db.CreateAccount(ctx, account)
	if err != nil {
		level.Error(logger).Log("err", err)
		return Account{}, err
	}

	logger.Log("create account", account.ID)

	return account, nil
}

// GetAccount API method for Service provide information about account.
func (s service) GetAccount(ctx context.Context, id int) (Account, error) {
	logger := log.With(s.logger, "method", "GetAccount")

	account, err := s.db.GetAccount(ctx, id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return Account{}, err
	}

	logger.Log("Get account", account.ID)

	return account, nil
}

// CreatePayment API method for Service provide create payment feature.
func (s service) CreatePayment(ctx context.Context, payment Payment) (Payment, error) {
	logger := log.With(s.logger, "method", "CreatePayment")

	payment, err := s.db.CreatePayment(ctx, payment)

	if err != nil {
		level.Error(logger).Log("err", err)
		return Payment{}, err
	}

	logger.Log("create payment", payment.ID)

	return payment, nil
}

// GetPayment API method for Service provide information about payment.
func (s service) GetPayment(ctx context.Context, id int) (Payment, error) {
	logger := log.With(s.logger, "method", "GetPayment")

	payment, err := s.db.GetPayment(ctx, id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return Payment{}, err
	}

	logger.Log("Get Payment", payment.ID)

	return payment, nil
}
