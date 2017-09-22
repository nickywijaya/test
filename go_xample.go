package go_xample

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type GoXample struct {
	database   DatabaseInterface
	messenger  MessengerInterface
	connection ConnectionInterface
}

type DatabaseInterface interface {
	InsertUser(context.Context, User) error
	FindUserByID(context.Context, int) (User, error)
	FindUserByCredential(context.Context, User) (User, error)
	FindInactiveUsers(context.Context) ([]User, error)
	InsertLoginHistory(context.Context, LoginHistory) error
	DeactivateUsers(context.Context, []User) error
}

type MessengerInterface interface {
	PublishLoginHistory(context.Context, LoginHistory) error
}

type ConnectionInterface interface {
	IsEmailValid(context.Context, string) (bool, error)
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name, omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

type LoginHistory struct {
	Username string    `json:"username"`
	LoginAt  time.Time `json:"login_at"`
}

func NewGoXample(db DatabaseInterface, msgr MessengerInterface, conn ConnectionInterface) GoXample {
	return GoXample{
		database:   db,
		messenger:  msgr,
		connection: conn,
	}
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

	err = g.database.InsertUser(ctx, user)
	if err != nil {
		return User{}, err
	}

	valid, err := g.connection.IsEmailValid(ctx, user.Email)
	if err != nil {
		return User{}, err
	}
	if !valid {
		return User{}, errors.New("Email Invalid")
	}

	return user, nil
}

func (g *GoXample) GetUserByID(ctx context.Context, id int) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, errors.New("Timeout")
	default:
	}

	user, err := g.database.FindUserByID(ctx, id)
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

	user, err = g.database.FindUserByCredential(ctx, user)
	if err != nil {
		return User{}, err
	}

	g.updateLoginHistory(ctx, user)

	return user, nil
}

func (g *GoXample) updateLoginHistory(ctx context.Context, user User) error {
	loginHistory := LoginHistory{
		Username: user.Username,
		LoginAt:  time.Now(),
	}

	return g.messenger.PublishLoginHistory(ctx, loginHistory)
}

func (g *GoXample) SaveLoginHistory(ctx context.Context, loginHistory LoginHistory) error {
	return g.database.InsertLoginHistory(ctx, loginHistory)
}

func (g *GoXample) DeactivateInactiveUsers(ctx context.Context) error {
	users, err := g.database.FindInactiveUsers(ctx)
	if err != nil {
		return err
	}

	err = g.database.DeactivateUsers(ctx, users)
	if err != nil {
		return err
	}

	return nil
}
