package handler

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ilya-korotya/solid/server"
	"github.com/ilya-korotya/solid/usecase"

	"github.com/ilya-korotya/solid/mock"

	"github.com/ilya-korotya/solid/entries"
)

func TestUsersHandler(t *testing.T) {
	type bind struct {
		w           *httptest.ResponseRecorder
		userUsecase *mock.UserUsecase
		code        int
		result      string
		err         error
	}
	mocks := map[string]bind{
		"Cannot return users from storage:": bind{
			w: httptest.NewRecorder(),
			userUsecase: &mock.UserUsecase{
				UsersFn: func() ([]*entries.User, error) {
					// TODO: return from database normal error
					return nil, usecase.InternalError.New("invalid connect to storage")
				},
			},
			code:   http.StatusInternalServerError,
			result: `{"error":"invalid connect to storage"}`,
			err:    usecase.InternalError.New("invalid connect to storage"),
		},
		"No users in the storage:": bind{
			w: httptest.NewRecorder(),
			userUsecase: &mock.UserUsecase{
				UsersFn: func() ([]*entries.User, error) {
					return nil, usecase.NotFound.New("cannot find users in storage")
				},
			},
			code:   http.StatusNotFound,
			result: `{"error":"cannot find users in storage"}`,
			err:    usecase.NotFound.New("cannot find users in storage"),
		},
		"Success return users:": bind{
			w: httptest.NewRecorder(),
			userUsecase: &mock.UserUsecase{
				UsersFn: func() ([]*entries.User, error) {
					return []*entries.User{
						&entries.User{
							FirstName:  "Lisa",
							SecondName: "Ann",
							Age:        45,
						},
						&entries.User{
							FirstName:  "Kendra",
							SecondName: "Lust",
							Age:        39,
						},
					}, nil
				},
			},
			code: http.StatusOK,
			// TODO: may be better make function for unmarshal json and make compare response and expected object
			// because always writing JSON is difficult and boring
			result: `[{"first_name":"Lisa","second_name":"Ann","age":45},{"first_name":"Kendra","second_name":"Lust","age":39}]`,
		},
	}
	for name, mock := range mocks {
		t.Run(name, func(t *testing.T) {
			c := server.NewContext(mock.w, nil)
			c.UserUsecase = mock.userUsecase
			err := users(c)
			if !mock.userUsecase.UsersInvoked {
				t.Errorf("Mock function users not called")
			}
			if mock.w.Code != mock.code {
				t.Errorf("HTTP status code from handler 'users' not match to expected code: %v != %v", mock.w.Code, mock.code)
			}
			if reflect.TypeOf(err) != reflect.TypeOf(mock.err) {
				t.Errorf("Not expected error type from handler 'users': %T != %T", err, mock.err)
			}
			r := mock.w.Body.String()
			if r != mock.result {
				t.Errorf("The handler 'users' response does not match the expected result: %s != %s", r, mock.result)
			}
		})
	}
}
