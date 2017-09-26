package database

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"

	gx "github.com/bukalapak/go-xample"
)

type MySQLMock struct{}

func (m *MySQLMock) InsertUser(ctx context.Context, user gx.User) error {
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	if user.Username == "bad-user" {
		return errors.New("Error! Bad user!")
	}

	return nil
}

func (m *MySQLMock) FindUserByID(ctx context.Context, id int) (gx.User, error) {
	select {
	case <-ctx.Done():
		return gx.User{}, errors.New("Timeout")
	default:
	}

	if id <= 0 {
		return gx.User{}, errors.New("Error! Bad user!")
	}

	user := unmarshalUser("user1.json")

	return user, nil
}

func (m *MySQLMock) FindUserByCredential(ctx context.Context, cred gx.User) (gx.User, error) {
	select {
	case <-ctx.Done():
		return gx.User{}, errors.New("Timeout")
	default:
	}

	if cred.Username == "user1" && cred.Password == "user1" {
		user := unmarshalUser("user1.json")
		return user, nil
	}

	return gx.User{}, errors.New("Error! Bad user!")
}

func (m *MySQLMock) InsertLoginHistory(ctx context.Context, loginHistory gx.LoginHistory) error {
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	if loginHistory.Username != "user1" {
		return errors.New("Error! Bad user!")
	}

	return nil
}

func (m *MySQLMock) FindInactiveUsers(ctx context.Context) ([]gx.User, error) {
	select {
	case <-ctx.Done():
		return []gx.User{}, errors.New("Timeout")
	default:
	}

	var users []gx.User

	if err, ok := ctx.Value("Error").(bool); ok && err {
		return users, errors.New("Error! Database crumbled!")
	}

	if count, ok := ctx.Value("Users").(int); ok && count <= 0 {
		return users, nil
	}

	user := unmarshalUser("user1.json")
	users = append(users, user)

	return users, nil
}

func (m *MySQLMock) DeactivateUsers(ctx context.Context, users []gx.User) error {
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	if len(users) == 0 {
		return errors.New("Error! Bad users!")
	}

	return nil
}

func readFixture(key string) ([]byte, error) {
	fxPath, _ := filepath.Abs(path.Join("../testdata", key))
	return ioutil.ReadFile(fxPath)
}

func unmarshalUser(key string) gx.User {
	var user gx.User

	b, _ := readFixture(key)
	json.Unmarshal(b, &user)

	return user
}
