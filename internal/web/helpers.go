package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func checkHeader(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("invalid content type header")
	}

	return nil
}

func getValidatorErrMsg(err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		e := errs[0]
		message := fmt.Sprintf("invalid field: %v", e.Field())
		return errors.New(message)
	}
	return err
}
