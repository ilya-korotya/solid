package server

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/ilya-korotya/solid/usecase"
)

type personaMock struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

func TestBind(t *testing.T) {
	type bind struct {
		r      *http.Request
		target *personaMock
		result *personaMock
		err    error
	}
	mocks := map[string]bind{
		"Invalid JSON in request": bind{
			r:   httptest.NewRequest("GET", "/", strings.NewReader(`{"name": "John" "age": 20}`)), // notice the ',' between 'name' and 'age'
			err: &json.SyntaxError{},
		},
		"Parameter which is not a link": bind{
			r:   httptest.NewRequest("GET", "/", strings.NewReader(`{"name": "John", "age": 20}`)),
			err: &json.InvalidUnmarshalError{},
		},
		"Query conversion success": bind{
			r:      httptest.NewRequest("GET", "/", strings.NewReader(`{"name": "John", "age": 20}`)),
			target: &personaMock{},
			result: &personaMock{
				Name: "John",
				Age:  20,
			},
		},
	}
	for name, mock := range mocks {
		t.Run(name, func(t *testing.T) {
			c := Context{
				r: mock.r,
			}
			err := c.Bind(mock.target)
			if reflect.TypeOf(err) != reflect.TypeOf(mock.err) {
				t.Errorf("Type errors after unmarshal JSON do not match: %T != %T", err, mock.err)
			}
			if !reflect.DeepEqual(mock.target, mock.result) {
				t.Errorf("Target is not equal to known: %+v != %+v", mock.target, mock.result)
			}
		})
	}
}

func TestResponse(t *testing.T) {
	type bind struct {
		w      *httptest.ResponseRecorder
		body   interface{}
		result string
		code   int
		err    error
	}
	mocks := map[string]bind{
		"Invalid JSON in response": bind{
			w:    httptest.NewRecorder(),
			body: math.Inf(1), // lol
			code: http.StatusInternalServerError,
			err:  &json.UnsupportedValueError{},
		},
		"Successful response to the client": bind{
			w: httptest.NewRecorder(),
			body: &personaMock{
				Name: "John",
				Age:  20,
			},
			result: `{"name":"John","age":20}`,
			code:   http.StatusOK,
		},
	}
	for name, mock := range mocks {
		t.Run(name, func(t *testing.T) {
			c := Context{
				w: mock.w,
			}
			err := c.Response(mock.code, mock.body)
			if reflect.TypeOf(err) != reflect.TypeOf(mock.err) {
				t.Errorf("Type errors after marshal JSON do not match: %T != %T", err, mock.err)
			}
			ct := c.w.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Invalid heading reflecting content type: %s != %s", ct, "application/json")
			}
			if mock.w.Code != mock.code {
				t.Errorf("Status Code changed in response: %v != %v", mock.w.Code, mock.code)
			}
			r := mock.w.Body.String()
			t.Log("Actual response body:", r)
			t.Log("Expected response body:", mock.result)
			if r != mock.result {
				t.Errorf("The answer to the client does not match: %s != %s", r, mock.result)
			}
		})
	}
}

func TestProcessError(t *testing.T) {
	type bind struct {
		w      *httptest.ResponseRecorder
		body   error
		result string
		code   int
		err    error
	}
	mocks := map[string]bind{
		"Success proccese 'NotFound' error:": bind{
			w:      httptest.NewRecorder(),
			body:   usecase.NotFound.New("user id=20 not found in storage"),
			result: `{"error":"user id=20 not found in storage"}`,
			code:   http.StatusNotFound,
			err:    usecase.NotFound.New("user id=20 not found in storage"),
		},
	}
	for name, mock := range mocks {
		t.Run(name, func(t *testing.T) {
			c := Context{
				w: mock.w,
			}
			err := c.ProcessError(mock.body)
			if reflect.TypeOf(err) != reflect.TypeOf(mock.err) {
				t.Errorf("Incorrectly proxied error type: %T != %T", err, mock.err)
			}
			ct := c.w.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Invalid heading reflecting content type: %s != %s", ct, "application/json")
			}
			if mock.w.Code != mock.code {
				t.Errorf("Incorrectly installed HTTP status code in response: %v != %v", mock.w.Code, mock.code)
			}
			r := mock.w.Body.String()
			t.Log("Actual response with error:", r)
			t.Log("Expected response with error:", mock.result)
			if r != mock.result {
				t.Errorf("The answer to the client does not match: %s != %s", r, mock.result)
			}
		})
	}
}
