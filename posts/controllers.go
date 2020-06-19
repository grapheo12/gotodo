package posts

import (
	"encoding/json"
	"gotodo/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //For SQLite Dialect
)

//Register the URLs
func Register(r *mux.Router) {
	r.HandleFunc("/", auth.LoginRequired(getAll)).Methods("GET")
	r.HandleFunc("/", auth.LoginRequired(savePost)).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", auth.LoginRequired(getOne)).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", auth.LoginRequired(delOne)).Methods("DELETE")
}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(auth.CtxUserString("user")).(auth.User)

	dbStr := "test.db"
	db, err := gorm.Open("sqlite3", dbStr)
	if err != nil {
		panic("Failed to connect to the Database!")
	}
	defer db.Close()

	result := (&Post{}).All(db, user)

	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	return
}

func savePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dbStr := "test.db"
	db, err := gorm.Open("sqlite3", dbStr)
	if err != nil {
		panic("Failed to connect to the Database!")
	}
	defer db.Close()

	var post Post

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if post.Deadline.Sub(time.Now()) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := r.Context().Value(auth.CtxUserString("user")).(auth.User)
	post.User = user
	post.Save(db)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "All OK"}`))
}

func getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	w.Header().Set("Content-Type", "application/json")
	dbStr := "test.db"
	db, err := gorm.Open("sqlite3", dbStr)
	if err != nil {
		panic("Failed to connect to the Database!")
	}
	defer db.Close()

	result := (&Post{}).One(db, id)

	if result.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := r.Context().Value(auth.CtxUserString("user")).(auth.User)
	if result.UserID != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func delOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	w.Header().Set("Content-Type", "application/json")
	dbStr := "test.db"
	db, err := gorm.Open("sqlite3", dbStr)
	if err != nil {
		panic("Failed to connect to the Database!")
	}
	defer db.Close()

	result := (&Post{}).One(db, id)

	if result.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user := r.Context().Value(auth.CtxUserString("user")).(auth.User)
	if result.UserID != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	result.Delete(db)

	w.WriteHeader(http.StatusFound)
	w.Write([]byte(`{"message": "Post deleted"}`))
}
