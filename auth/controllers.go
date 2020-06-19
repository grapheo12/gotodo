package auth

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //For SQLite Dialect
)

//Claims jwt claim
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//Register auth urls
func Register(r *mux.Router) {
	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/refresh", refresh).Methods("GET")
}

func signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data map[string]string
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok1 := data["username"]
	_, ok2 := data["password"]

	if !(ok1 && ok2) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbStr := "test.db"
	db, err := gorm.Open("sqlite3", dbStr)
	if err != nil {
		panic("Failed to connect to the Database!")
	}
	defer db.Close()

	checkUser := (&User{}).One(db, data["username"])
	if checkUser.ID != 0 {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte(`{"message": "User already exists"}`))
		return
	}

	var user User
	user.Username = data["username"]

	hash := md5.Sum([]byte(data["password"]))
	user.PwdHash = base64.StdEncoding.EncodeToString(hash[:])

	user.Save(db)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "User Created"}`))

}

func login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data map[string]string
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok1 := data["username"]
	_, ok2 := data["password"]

	if !(ok1 && ok2) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbStr := "test.db"
	db, err := gorm.Open("sqlite3", dbStr)
	if err != nil {
		panic("Failed to connect to the Database!")
	}
	defer db.Close()

	checkUser := (&User{}).One(db, data["username"])

	hash := md5.Sum([]byte(data["password"]))
	if base64.StdEncoding.EncodeToString(hash[:]) != checkUser.PwdHash {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Wrong Username/Password"}`))
		return
	}

	jwtKey := []byte("OneRandomSecretKey!!@@!")

	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &Claims{
		Username: checkUser.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{"token": "` + tokenStr + `"}`))
}

//Token must come from the Bearer header
func refresh(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Bearer")
	if tokenStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jwtKey := []byte("OneRandomSecretKey!!@@!")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(30 * time.Second)
	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenStr, err = token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{"token": "` + tokenStr + `"}`))

}
