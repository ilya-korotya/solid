package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/lib/pq"

	"github.com/ilya-korotya/solid/server"
	"github.com/ilya-korotya/solid/usecase"

	"github.com/ilya-korotya/solid/mock"
)

func TestUserCreate(t *testing.T) {
	type bind struct {
		w               *httptest.ResponseRecorder
		r               *http.Request
		userUsecase     *mock.UserUsecase
		registerInvoked mock.CallMock // check call RegisterFn
		code            int
		result          string
		err             error
	}
	mocks := map[string]bind{
		"Invalid bind JSON from request": bind{
			w:           httptest.NewRecorder(),
			r:           httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name", John}`)), // invalid JSON. Value 'name' wrote withot ""
			userUsecase: &mock.UserUsecase{},
			code:        http.StatusBadRequest,
			result:      `{"error":"invalid character ',' after object key"}`, // TODO: is it worth it?
			err:         &json.SyntaxError{},
		},
		"Cannot register new user. Invalid connect to storage:": bind{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "John"}`)),
			userUsecase: &mock.UserUsecase{
				RegisterFn: func(c *usecase.Client) (bool, error) {
					return false, pq.ErrChannelNotOpen
				},
			},
			registerInvoked: true,
			code:            http.StatusInternalServerError,
			result:          `{"error":"pq: channel is not open"}`,
			err:             pq.ErrChannelNotOpen,
		},
		"Success create new user": bind{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"first_name":"Lisa", "second_name":"Ann", "age":20}`)),
			userUsecase: &mock.UserUsecase{
				RegisterFn: func(c *usecase.Client) (bool, error) {
					return true, nil
				},
			},
			registerInvoked: true,
			code:            http.StatusOK,
			result:          `{"first_name":"Lisa","second_name":"Ann","age":20,"tits_size":0}`,
		},
	}
	for name, mock := range mocks {
		c := server.NewContext(mock.w, mock.r)
		c.UserUsecase = mock.userUsecase
		err := userCreate(c)
		t.Run(name, func(t *testing.T) {
			if mock.userUsecase.RegisterInvoked != mock.registerInvoked {
				t.Errorf("Callback 'Register' have status '%s'. When it must have status: '%s'", mock.userUsecase.RegisterInvoked, mock.registerInvoked)
			}
			if reflect.TypeOf(err) != reflect.TypeOf(mock.err) {
				t.Errorf("Handler 'user_create' error does not match expected error: %T != %T", err, mock.err)
			}
			if mock.w.Code != mock.code {
				t.Errorf("Status code response from 'user_create' does not match the expected: %v != %v", mock.w.Code, mock.code)
			}
			r := mock.w.Body.String()
			if r != mock.result {
				t.Errorf("The message sent to 'user_create' handlers does not match the expected: %s != %s", r, mock.result)
			}
		})
	}
}
