package go_xample

import (
	"encoding/json"
)

type GoXample struct {
	// DB
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
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
