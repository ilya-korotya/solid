package entries

type User struct {
	Name string
	Age  int
}

type UserStore interface {
	Get(id int) *User
	Create(*User) *User
}
