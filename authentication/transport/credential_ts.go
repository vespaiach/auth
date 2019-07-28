package transport

import (
	"context"
	"encoding/json"
	"net/http"
)

func decodeVerifyLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	if len(body.Password) == 0 || len(body.Username) == 0 {
		return nil, ErrInvalidPayload
	}

	return verifyLoginRequest{
		body.Username,
		body.Password,
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e.Error() != "" {
		encodeError(ctx, e, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
