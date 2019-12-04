// server.go provide Service interface, HTTP server and middleware.
package account

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// Service interface.
type Service interface {
	CreateAccount(ctx context.Context, account Account) (Account, error)
	GetAccount(ctx context.Context, id int) (Account, error)

	CreatePayment(ctx context.Context, payment Payment) (Payment, error)
	GetPayment(ctx context.Context, id int) (Payment, error)
}

// NewHTTPServer - setups application API URI router.
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/accounts").Handler(httptransport.NewServer(
		endpoints.CreateAccount,
		decodeCreateAccountReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/accounts/{id}").Handler(httptransport.NewServer(
		endpoints.GetAccount,
		decodeIDReq,
		encodeResponse,
	))

	r.Methods("POST").Path("/payments").Handler(httptransport.NewServer(
		endpoints.CreatePayment,
		decodePaymentReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/payments/{id}").Handler(httptransport.NewServer(
		endpoints.GetPayment,
		decodeIDReq,
		encodeResponse,
	))

	return r

}

// commonMiddleware - add HTTP header for JSON format.
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
