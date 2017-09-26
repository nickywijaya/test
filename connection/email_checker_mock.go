package connection

import (
	"context"
	"errors"
)

type EmailCheckerMock struct{}

func (e *EmailCheckerMock) IsEmailValid(ctx context.Context, email string) (bool, error) {
	if email == "bad@email.com" {
		return false, errors.New("Error! Bad email!")
	} else if email == "invalid@email.com" {
		return false, nil
	}

	return true, nil
}
