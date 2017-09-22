package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	gx "github.com/bukalapak/go-xample"
)

type MySQL struct {
	db *sql.DB
}

type Option struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Charset  string
}

func NewMySQL(opt Option) (*MySQL, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", opt.User, opt.Password, opt.Host, opt.Port, opt.Database, opt.Charset))
	if err != nil {
		return &MySQL{}, err
	}

	return &MySQL{db: db}, nil
}

func (m *MySQL) InsertUser(ctx context.Context, user gx.User) error {
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	_, err := m.db.Exec("INSERT INTO users(username, name, password, email) VALUES(?, ?, ?, ?)", user.Username, user.Name, user.Password, user.Email)
	return err
}

func (m *MySQL) FindUserByID(ctx context.Context, id int) (gx.User, error) {
	select {
	case <-ctx.Done():
		return gx.User{}, errors.New("Timeout")
	default:
	}

	var user gx.User

	err := m.db.QueryRow("SELECT id, name, username, password, email, active FROM users WHERE active = true AND id = ?", id).Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email, &user.Active)
	if err != nil {
		return gx.User{}, err
	}

	return user, nil
}

func (m *MySQL) FindUserByCredential(ctx context.Context, cred gx.User) (gx.User, error) {
	select {
	case <-ctx.Done():
		return gx.User{}, errors.New("Timeout")
	default:
	}

	var user gx.User

	err := m.db.QueryRow("SELECT id, name, username, password, email, active FROM users WHERE active = true AND username = ? AND password = ?", cred.Username, cred.Password).Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email, &user.Active)
	if err != nil {
		return gx.User{}, err
	}

	return user, nil
}

func (m *MySQL) InsertLoginHistory(ctx context.Context, loginHistory gx.LoginHistory) error {
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	_, err := m.db.Exec("INSERT INTO login_histories (username, login_at) VALUES(?, ?)", loginHistory.Username, loginHistory.LoginAt)
	return err
}

func (m *MySQL) FindInactiveUsers(ctx context.Context) ([]gx.User, error) {
	select {
	case <-ctx.Done():
		return []gx.User{}, errors.New("Timeout")
	default:
	}

	var users []gx.User

	rows, err := m.db.Query(`
			SELECT users.username
			FROM users
			INNER JOIN login_histories
				ON users.username = login_histories.username
			WHERE users.active = true
			GROUP BY users.username
			HAVING DATEDIFF(?, MAX(login_histories.login_at)) >= 30
	`, time.Now())
	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		var user gx.User

		if err = rows.Scan(&user.Username); err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (m *MySQL) DeactivateUsers(ctx context.Context, users []gx.User) error {
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	var usernames []interface{}

	for _, user := range users {
		usernames = append(usernames, user.Username)
	}

	_, err := m.db.Exec(`
		UPDATE users
		SET active = false
		WHERE username IN (? `+strings.Repeat(",?", len(usernames)-1)+`)
	`, usernames...)
	if err != nil {
		return err
	}

	return nil
}
