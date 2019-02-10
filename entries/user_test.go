package entries_test

import (
	"testing"

	"github.com/ilya-korotya/solid/entries"
)

type test struct {
	firstName  string
	secondName string
	age        uint8
	err        error
}

func TestUserConstructor(t *testing.T) {
	mocks := map[string]test{
		"Aktrisa is too old for filming": test{
			firstName:  "Lisa",
			secondName: "Ann",
			age:        35,
			err:        entries.ErrorOldFag,
		},
		"Aktrisa is too young for filming": test{
			firstName:  "Alex",
			secondName: "Texas",
			age:        17,
			err:        entries.ErrorNewFag,
		},
		"The actress is perfect for filming": test{
			firstName:  "Kendra",
			secondName: "Lust",
			age:        18,
		},
	}
	for name, mock := range mocks {
		t.Run(name, func(t *testing.T) {
			if _, err := entries.NewUser(mock.firstName, mock.secondName, mock.age); err != mock.err {
				t.Errorf("Type of creation errors do not match: %s != %s", err, mock.err)
			}
		})
	}
}
