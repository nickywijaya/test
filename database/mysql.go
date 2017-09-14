package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func NewMySQL(opt Option) (MySQL, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", opt.User, opt.Password, opt.Host, opt.Port, opt.Database, opt.Charset))
	if err != nil {
		return MySQL{}, err
	}

	return MySQL{db: db}, nil
}

func (m MySQL) FindUserByID(ctx context.Context, id int) (gx.User, error) {
	var user gx.User

	err := m.db.QueryRow("SELECT * FROM users WHERE id=?").Scan(&user)
	if err != nil {
		return gx.User{}, err
	}

	return user, nil
}

func (m MySQL) InsertUser(ctx context.Context, user gx.User) error {
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	_, err := m.db.Exec("INSERT INTO users(username, name, password) VALUES(?, ?, ?)", user.Name, user.Username, user.Password)
	return err
}
