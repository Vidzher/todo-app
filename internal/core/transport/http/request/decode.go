package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("error decoding json: %w", err)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
