// endpoint.go service endpoints
package account

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

// CreateAccount - provide create account feature
// GetAccount    - provide information about account
// CreatePayment - provide create payment feature
// GetPayment    - provide information about payment
type Endpoints struct {
	CreateAccount endpoint.Endpoint
	GetAccount    endpoint.Endpoint
	CreatePayment endpoint.Endpoint
	GetPayment    endpoint.Endpoint
}

// MakeEndpoints creates application endpoints.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateAccount: makeCreateAccountEndpoint(s),
		GetAccount:    makeGetAccountEndpoint(s),
		CreatePayment: makeCreatePaymentEndpoint(s),
		GetPayment:    makeGetPaymentEndpoint(s),
	}
}

// makeCreateAccountEndpoint creates CreateAccount application endpoint.
func makeCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// validate Account schema
		req := request.(Account)
		return s.CreateAccount(ctx, req)
	}
}

// makeGetAccountEndpoint creates GetAccount application endpoint.
func makeGetAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// validate GetIDRequest schema
		req := request.(GetIDRequest)
		return s.GetAccount(ctx, req.ID)
	}
}

// makeCreatePaymentEndpoint creates CreatePayment application endpoint.
func makeCreatePaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// validate Payment schema
		req := request.(Payment)
		return s.CreatePayment(ctx, req)
	}
}

// makeGetPaymentEndpoint creates GetPayment application endpoint.
func makeGetPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// validate GetIDRequest schema
		req := request.(GetIDRequest)
		return s.GetPayment(ctx, req.ID)
	}
}
