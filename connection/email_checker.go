// Package connection contains implementation of connection service.
// Any connection to third party service should be implemented here.
package connection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rubyist/circuitbreaker"
)

// EmailChecker holds the functionality to call Email Checker service.
type EmailChecker struct {
	cb     *circuit.Breaker
	client *http.Client
	option Option
}

// Option holds all necessary options to make connection.
type Option struct {
	Timeout time.Duration
}

// Email holds the response from Email Checker service.
type Email struct {
	Valid bool `json:"valid"`
}

// NewEmailChecker returns a pointer of EmailChecker instance
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

// IsEmailValid checks whether an email is valid.
func (e *EmailChecker) IsEmailValid(ctx context.Context, email string) (bool, error) {
	if e.cb.Ready() {
		request, _ := http.NewRequest("GET", fmt.Sprintf("%s?email=%s", os.Getenv("EMAIL_CHECKER_URL"), email), nil)

		reqID, _ := ctx.Value("Request-ID").(string)
		actor, _ := ctx.Value("Authorization").(string)

		request.Header.Set("Accept", "application/json")
		request.Header.Set("Request-ID", reqID)
		request.Header.Set("Authorization", actor)

		var (
			err      error
			response *http.Response
		)

		retry := 0
		for retry < 3 {
			request.Header.Set("Retry", strconv.Itoa(retry))

			response, err = e.client.Do(request)
			if (err != nil || response.StatusCode != 200) && retry < 3 {
				retry++
				continue
			}

			if err != nil {
				e.cb.Fail()
				return false, err
			}

			if response.StatusCode != 200 {
				e.cb.Fail()
				return false, errors.New(fmt.Sprintf("HTTP Response Code: %d", response.StatusCode))
			}

			break
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			e.cb.Success()
			return false, err
		}

		var mail Email

		err = json.Unmarshal(body, &mail)
		if err != nil {
			e.cb.Success()
			return false, err
		}

		e.cb.Success()
		return mail.Valid, nil
	} else {
		return false, errors.New("Circuit Breaker is in Trip State")
	}
}
