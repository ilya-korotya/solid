package entries

import "errors"

const (
	startIdealAge = 18
	endIdealAge   = 25 // lol
)

// Basic error for user
var (
	ErrorOldFag = errors.New("you are very old for this shit")
	ErrorNewFag = errors.New("you are very young for this shit")
)

type UserStore interface {
	User(id string) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) error
	DeleteUser(id string) error
}

type User struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Age        uint8  `json:"age"`
}

func NewUser(firstName, secondName string, age uint8) (*User, error) {
	if age < startIdealAge {
		return nil, ErrorNewFag
	}
	if age > endIdealAge {
		return nil, ErrorOldFag
	}
	return &User{
		FirstName:  firstName,
		SecondName: secondName,
		Age:        age,
	}, nil
}
