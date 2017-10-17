// Package go_xample is the main package for GoXample project.
// It contains all definitions and implementation of the project.
// It also specifies its dependencies.
package go_xample

import (
	"context"
	"errors"
	"time"
)

// GoXample is the main struct that holds all business logic.
// Its functions explain its service and purpose.
type GoXample struct {
	database   Database
	messenger  Messenger
	connection Connection
}

// Database is a contract for database client.
// The functions are specific.
// They depends on GoXample's needs.
type Database interface {
	InsertUser(context.Context, User) error
	FindUserByID(context.Context, int) (User, error)
	FindUserByCredential(context.Context, User) (User, error)
	FindInactiveUsers(context.Context) ([]User, error)
	InsertLoginHistory(context.Context, LoginHistory) error
	DeactivateUsers(context.Context, []User) error
}

// Messenger is a contract for messenger client.
// The functions are specific.
// They depends on GoXample's needs.
type Messenger interface {
	PublishLoginHistory(context.Context, LoginHistory) error
}

// Connection is a contract for third party service.
// The functions are specific.
// They depends on GoXample's needs.
type Connection interface {
	IsEmailValid(context.Context, string) (bool, error)
}

// User holds data for user.
// It has six fields.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name, omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

// LoginHistory holds data for user's login history.
// It only has two fields: Username and LoginAt.
type LoginHistory struct {
	Username string    `json:"username"`
	LoginAt  time.Time `json:"login_at"`
}

// NewGoXample returns a pointer of GoXample instance.
// It takes three parameters.
func NewGoXample(db Database, msgr Messenger, conn Connection) *GoXample {
	return &GoXample{
		database:   db,
		messenger:  msgr,
		connection: conn,
	}
}

// CreateUser is a function to create user and write it to database
func (g *GoXample) CreateUser(ctx context.Context, user User) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, errors.New("Timeout")
	default:
	}

	err := g.database.InsertUser(ctx, user)
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

// GetUserByID is a function to retrieve user data by their id.
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

// GetUserByCredential is a function to retrieve user data by their username and password.
func (g *GoXample) GetUserByCredential(ctx context.Context, user User) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, errors.New("Timeout")
	default:
	}

	user, err := g.database.FindUserByCredential(ctx, user)
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

// SaveLoginHistory is a function to save LoginHistory to database.
func (g *GoXample) SaveLoginHistory(ctx context.Context, loginHistory LoginHistory) error {
	return g.database.InsertLoginHistory(ctx, loginHistory)
}

// DeactivateInactiveUsers is a function to turn the value of field `active` from true to false.
// The affected users are users who haven't done any login activity in the last 30 days.
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
