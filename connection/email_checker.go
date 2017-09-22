package connection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rubyist/circuitbreaker"
)

type EmailChecker struct {
	cb     *circuit.Breaker
	client *http.Client
	option Option
}

type Option struct {
	Timeout time.Duration
}

type Email struct {
	Valid bool `json:"valid"`
}

func NewEmailChecker(opt Option) *EmailChecker {
	if opt.Timeout == 0 {
		opt.Timeout = 3 * time.Second
	}

	client := &http.Client{
		Timeout: opt.Timeout,
	}

	// create a circuit breaker that will trip if there are at least 10 errors out of 100 requests (10%)
	cb := circuit.NewRateBreaker(0.1, 100)

	return &EmailChecker{
		cb:     cb,
		client: client,
		option: opt,
	}
}

func (a *EmailChecker) IsEmailValid(ctx context.Context, email string) (bool, error) {
	if a.cb.Ready() {
		request, _ := http.NewRequest("GET", fmt.Sprintf("%s?email=%s", os.Getenv("EMAIL_CHECKER_URL"), email), nil)

		reqID, _ := ctx.Value("Request-ID").(string)

		request.Header.Set("Accept", "application/json")
		request.Header.Set("Request-ID", reqID)

		var (
			err      error
			response *http.Response
		)

		retry := 0
		for retry < 3 {
			response, err = a.client.Do(request)
			if (err != nil || response.StatusCode != 200) && retry < 3 {
				retry++
				continue
			}

			if err != nil {
				a.cb.Fail()
				return false, err
			}

			if response.StatusCode != 200 {
				a.cb.Fail()
				return false, errors.New(fmt.Sprintf("HTTP Response Code: %d", response.StatusCode))
			}
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			a.cb.Success()
			return false, err
		}

		var mail Email

		err = json.Unmarshal(body, &mail)
		if err != nil {
			a.cb.Success()
			return false, err
		}

		a.cb.Success()
		return mail.Valid, nil
	} else {
		return false, errors.New("Circuit Breaker is in Trip State")
	}
}
