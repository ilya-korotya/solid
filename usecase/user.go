package usecase

import (
	"github.com/ilya-korotya/solid/entries"
)

// TODO: think about the proper naming of variables and other things

type UserUsecase interface {
	Register(*Client) (bool, error)
	Users() ([]*entries.User, error)
}

type UserInteractor struct {
	UserStore entries.UserStore
}

func NewUserInteractor(us entries.UserStore) *UserInteractor {
	return &UserInteractor{
		UserStore: us,
	}
}

type Client struct {
	entries.User
	TitsSize uint8 `json:"tits_size"`
}

func (u *UserInteractor) Register(client *Client) (bool, error) {
	user, err := entries.NewUser(client.FirstName, client.SecondName, client.Age)
	if err != nil {
		return false, BadRequest.FromError(err)
	}
	err = u.UserStore.CreateUser(user)
	return true, InternalError.FromError(err)
}

func (u *UserInteractor) Users() ([]*entries.User, error) {
	return u.UserStore.Users()
}

// Ideally, we should have created a separate response structure.
