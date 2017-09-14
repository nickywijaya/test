package go_xample

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type GoXample struct {
	db DBInterface
}

type DBInterface interface {
	InsertUser(context.Context, User) error
	FindUserByID(context.Context, int) (User, error)
	FindUserByCredential(context.Context, User) (User, error)
	InsertLoginHistory(context.Context, User, time.Time) error
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name, omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

func NewGoXample(db DBInterface) GoXample {
	return GoXample{db: db}
}

func (g *GoXample) CreateUser(ctx context.Context, data string) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, errors.New("Timeout")
	default:
	}

	var user User

	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return User{}, err
	}

	err = g.db.InsertUser(ctx, user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (g *GoXample) GetUserByID(ctx context.Context, id int) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, errors.New("Timeout")
	default:
	}

	user, err := g.db.FindUserByID(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (g *GoXample) GetUserByCredential(ctx context.Context, data string) (User, error) {
	var user User

	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return User{}, err
	}

	user, err = g.db.FindUserByCredential(ctx, user)
	if err != nil {
		return User{}, err
	}

	go g.updateLoginHistory(ctx, user)

	return user, nil
}

func (g *GoXample) updateLoginHistory(ctx context.Context, user User) error {
	loginAt := time.Now()
	return g.db.InsertLoginHistory(ctx, user, loginAt)
}
