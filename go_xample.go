package go_xample

import (
	"encoding/json"
)

type GoXample struct {
	// DB
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name, omitempty"`
	Password string `json:"password, omitempty"`
}

func NewGoXample() GoXample {
	return GoXample{}
}

func (g *GoXample) CreateUser(data string) (User, error) {
	var user User

	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (g *GoXample) GetUser(id int) (User, error) {
	user := User{ID: id}
	return user, nil
}
