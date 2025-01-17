package mock

import (
	"github.com/ilya-korotya/solid/entries"
	"github.com/ilya-korotya/solid/usecase"
)

// CallMock display callback call status
type CallMock bool

func (c CallMock) String() string {
	if c {
		return "CALL"
	}
	return "NO CALL"
}

// UserStore implement mock for user database
type UserStore struct {
	UserFn            func(id string) (*entries.User, error)
	UserInvoked       CallMock
	UsersFn           func() ([]*entries.User, error)
	UsersInvoked      CallMock
	CreateUserFn      func(user *entries.User) error
	CreateUserInvoked CallMock
	DeleteUserFn      func(id string) error
	DeleteUserInvoked CallMock
}

func (u *UserStore) User(id string) (*entries.User, error) {
	u.UserInvoked = true
	return u.UserFn(id)
}

func (u *UserStore) Users() ([]*entries.User, error) {
	u.UsersInvoked = true
	return u.UsersFn()
}

func (u *UserStore) CreateUser(user *entries.User) error {
	u.CreateUserInvoked = true
	return u.CreateUserFn(user)
}

func (u *UserStore) DeleteUser(id string) error {
	u.DeleteUserInvoked = true
	return u.DeleteUserFn(id)
}

// UserUsecase implement mock for user usecase
type UserUsecase struct {
	RegisterFn      func(*usecase.Client) (bool, error)
	RegisterInvoked CallMock
	UsersFn         func() ([]*entries.User, error)
	UsersInvoked    CallMock
}

func (u *UserUsecase) Register(c *usecase.Client) (bool, error) {
	u.RegisterInvoked = true
	return u.RegisterFn(c)
}

func (u *UserUsecase) Users() ([]*entries.User, error) {
	u.UsersInvoked = true
	return u.UsersFn()
}
