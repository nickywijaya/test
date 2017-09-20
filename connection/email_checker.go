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
)

type EmailChecker struct {
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

	return &EmailChecker{
		client: client,
		option: opt,
	}
}

func (a *EmailChecker) IsEmailValid(ctx context.Context, email string) (bool, error) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s?email=%s", os.Getenv("EMAIL_CHECKER_URL"), email), nil)

	reqID, _ := ctx.Value("Request-ID").(string)

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Request-ID", reqID)

	response, err := a.client.Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode != 200 {
		return false, errors.New(fmt.Sprintf("HTTP Response Code: %d", response.StatusCode))
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	var mail Email

	err = json.Unmarshal(body, &mail)
	if err != nil {
		return false, err
	}

	return mail.Valid, nil
}
