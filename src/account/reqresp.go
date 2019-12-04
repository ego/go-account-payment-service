// reqresp.go provide struct for schema validation for Service API and
// encode and decode feature for request and response.
package account

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetIDRequest schema for URI parameter `id`.
type GetIDRequest struct {
	ID int `json:"id"`
}

// encodeResponse - encode API response.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// decodeCreateAccountReq - decode create account request.
func decodeCreateAccountReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req Account
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// decodeIDReq - gets URI parameter `id` and cast it.
func decodeIDReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetIDRequest
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	req = GetIDRequest{
		ID: id,
	}
	return req, nil
}

// decodePaymentReq - decode payment body request as JSON object.
func decodePaymentReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req Payment
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
