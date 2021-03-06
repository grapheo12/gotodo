package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //For SQLite Dialect
)

//CtxUserString type for using with context
type CtxUserString string

//LoginRequired Middleware to protect endpoints
func LoginRequired(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		ctx := r.Context()

		dbStr := "test.db"
		db, err := gorm.Open("sqlite3", dbStr)
		if err != nil {
			panic("Failed to connect to the Database!")
		}
		defer db.Close()

		user := (&User{}).One(db, claims.Username)
		newctx := context.WithValue(ctx, CtxUserString("user"), user)
		req := r.WithContext(newctx)

		next(w, req)
	}
}
