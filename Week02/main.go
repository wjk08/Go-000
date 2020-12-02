package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

func main() {
	api()
}

type User struct {
	name string
}

func QueryUser() (*User, error) {
	//error
	//return nil, errors.New("test")
	//return &User{"go"}, nil
	return nil, errors.Wrapf(sql.ErrNoRows, "error")
}

func service() (string, error) {
	user, err := QueryUser()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		user = &User{"default"}
	}

	return fmt.Sprintf("hello %s", user.name), nil
}

func api() {
	message, err := service()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Println(message)
}
