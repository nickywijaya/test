package messenger

import (
	"context"
	"errors"

	gx "github.com/bukalapak/go-xample"
)

type RabbitMQMock struct{}

func (r *RabbitMQMock) PublishLoginHistory(ctx context.Context, loginHistory gx.LoginHistory) error {
	if loginHistory.Username == "bad-user" {
		return errors.New("Error! Bad user!")
	}

	return nil
}
