package main

import (
	"gotodo/auth"
	"gotodo/posts"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //For SQLite Dialect
)

func main() {
	dbStr := "test.db"
	db, err := gorm.Open("sqlite3", dbStr)
	if err != nil {
		panic("Failed to connect to the Database!")
	}

	posts.AutoInit(db)
	auth.AutoInit(db)

	db.Close()

	router := mux.NewRouter()
	authSubRoute := router.PathPrefix("/auth/").Subrouter()
	postsSubRoute := router.PathPrefix("/posts/").Subrouter()
	posts.Register(postsSubRoute)
	auth.Register(authSubRoute)

	log.Fatal(http.ListenAndServe(":8080", router))
}
