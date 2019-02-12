package postgres

import (
	"database/sql"

	"github.com/ilya-korotya/solid/entries"
	_ "github.com/lib/pq"
)

const (
	getUser    = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	getUsers   = "SELECT first_name, second_name, age FROM users"
	createUser = "INSERT INTO users (first_name, second_name, age) VALUES ($1, $2, $3)"
	deleteUser = "DELETE FROM WHERE id=$1"
)

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

func (u *UserStore) User(id string) (*entries.User, error) {
	user := &entries.User{}
	err := u.DB.QueryRow(getUser, id).Scan(
		&user.FirstName,
		&user.SecondName,
		&user.Age,
	)
	return user, err
}

func (u *UserStore) Users() ([]*entries.User, error) {
	users := []*entries.User{}
	rows, err := u.DB.Query(getUsers)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := &entries.User{}
		if err := rows.Scan(
			&user.FirstName,
			&user.SecondName,
			&user.Age,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (u *UserStore) CreateUser(user *entries.User) error {
	_, err := u.DB.Exec(createUser, user.FirstName, user.SecondName, user.Age)
	return err
}

func (u *UserStore) DeleteUser(id string) error {
	_, err := u.DB.Exec(deleteUser, id)
	return err
}
