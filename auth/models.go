package auth

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //For SQLite Dialect
)

//User Model for auth
type User struct {
	gorm.Model
	Username string `json:"username"`
	PwdHash  string `json:"pwdhash"`
}

//AutoInit to Run Migrations
func AutoInit(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

//Save to register user
func (u *User) Save(db *gorm.DB) {
	db.Create(u)
}

//Delete particular post
func (u *User) Delete(db *gorm.DB) {
	db.Delete(u)
}

//One get one user
func (u *User) One(db *gorm.DB, username string) User {
	var user User

	db.First(&user, "username = ?", username)
	return user
}
