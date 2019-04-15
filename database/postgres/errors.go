package postgres

import "github.com/lib/pq"

// Custom errors for postgres database
var (
	ErrUniqUser = &pq.Error{
		Message: "user already exist in database",
		Detail:  "user already exist in database (inner)",
	}
)
