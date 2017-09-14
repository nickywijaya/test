package go_xample

import (
	"context"
	"encoding/json"
	"errors"
)

type GoXample struct {
	db DBInterface
}

type DBInterface interface {
	FindUserByID(context.Context, int) (User, error)
	InsertUser(context.Context, User) error
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name, omitempty"`
	Username string `json:"username, omitempty"`
	Password string `json:"password, omitempty"`
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

	// TODO: get from DB
	return user, nil
}
