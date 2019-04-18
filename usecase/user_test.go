package usecase_test

// we had to declare a package of this type because we have circular dependencies with mock package
// there is another solution. Break mock package into subpackages, which removes cyclic dependency

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ilya-korotya/solid/entries"
	"github.com/ilya-korotya/solid/mock"
	"github.com/ilya-korotya/solid/usecase"
)

func TestRegister(t *testing.T) {
	type bind struct {
		client *usecase.Client
		store  *mock.UserStore
		result bool
		err    error
	}
	mocks := map[string]bind{
		"Client registration failed. Age is too young:": bind{
			client: &usecase.Client{
				User: entries.User{
					FirstName:  "Sophie",
					SecondName: "Dee",
					Age:        8,
				},
				TitsSize: 3,
			},
			result: false,
			err:    usecase.BadRequest.Wrap(entries.ErrorNewFag),
		},
		"Client registration failed. Age is too old:": bind{
			client: &usecase.Client{
				User: entries.User{
					FirstName:  "Sophie",
					SecondName: "Dee",
					Age:        45,
				},
				TitsSize: 3,
			},
			result: false,
			err:    usecase.BadRequest.Wrap(entries.ErrorOldFag),
		},
		"Unable to write created user to storage": bind{
			client: &usecase.Client{
				User: entries.User{
					FirstName:  "Sophie",
					SecondName: "Dee",
					Age:        20,
				},
				TitsSize: 3,
			},
			store: &mock.UserStore{
				CreateUserFn: func(u *entries.User) error {
					return errors.New("invalid connect to store")
				},
			},
			result: false,
			err:    errors.New("invalid connect to store"),
		},
		"Success create new client:": bind{
			client: &usecase.Client{
				User: entries.User{
					FirstName:  "Sophie",
					SecondName: "Dee",
					Age:        20,
				},
				TitsSize: 3,
			},
			store: &mock.UserStore{
				CreateUserFn: func(u *entries.User) error {
					return nil
				},
			},
			result: true,
		},
	}
	for name, mock := range mocks {
		t.Run(name, func(t *testing.T) {
			u := usecase.NewUserInteractor(mock.store)
			result, err := u.Register(mock.client)
			if reflect.TypeOf(err) != reflect.TypeOf(mock.err) {
				t.Errorf("The type of error returned by the function is not equal to the expected type of error: %T != %T", err, mock.err)
			}
			if result != mock.result {
				t.Errorf("The result of the function is not equal to the expected: %v != %v", result, mock.result)
			}
		})
	}
}

func TestUsers(t *testing.T) {
	type bind struct {
		store  *mock.UserStore
		result []*entries.User
		err    error
	}
	mocks := map[string]bind{
		"Unable to get users back from storage:": bind{
			store: &mock.UserStore{
				UsersFn: func() ([]*entries.User, error) {
					return nil, usecase.InternalError.New("invalid connect to storage")
				},
			},
			err: usecase.InternalError.New("invalid connect to storage"),
		},
		"Success return users:": bind{
			store: &mock.UserStore{
				UsersFn: func() ([]*entries.User, error) {
					return []*entries.User{
						&entries.User{
							FirstName:  "Sophie",
							SecondName: "Dee",
							Age:        20,
						},
						&entries.User{
							FirstName:  "Kendra",
							SecondName: "Lust",
							Age:        22,
						},
					}, nil
				},
			},
			result: []*entries.User{
				&entries.User{
					FirstName:  "Sophie",
					SecondName: "Dee",
					Age:        20,
				},
				&entries.User{
					FirstName:  "Kendra",
					SecondName: "Lust",
					Age:        22,
				},
			},
		},
	}
	for name, mock := range mocks {
		t.Run(name, func(t *testing.T) {
			u := usecase.NewUserInteractor(mock.store)
			result, err := u.Users()
			if reflect.TypeOf(err) != reflect.TypeOf(mock.err) {
				t.Errorf("The type of error returned by the function is not equal to the expected type of error: %T != %T", err, mock.err)
			}
			if !reflect.DeepEqual(result, mock.result) {
				t.Errorf("The result of the function is not equal to the expected: %v != %v", result, mock.result)
			}
		})
	}
}
