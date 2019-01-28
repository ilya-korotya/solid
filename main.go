package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://root:@postgres:5432/solid")
	if err != nil {
		fmt.Println("Invalid open connect to database:", err)
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Invalid connect to database:", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})
	http.ListenAndServe(":8080", nil)
}
