package posts

import (
	"gotodo/auth"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //For SQLite Dialect
)

//Post Model to save todo posts
type Post struct {
	gorm.Model
	Head        string    `json:"head"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	UserID      uint      `json:"userId"`
	User        auth.User
}

//AutoInit to Run Migrations
func AutoInit(db *gorm.DB) {
	db.AutoMigrate(&Post{})
}

//Save to save Posts
func (p *Post) Save(db *gorm.DB) {
	db.Create(p)
}

//All Get all posts
func (p *Post) All(db *gorm.DB, user auth.User) []Post {
	var resultset []Post
	db.Model(&user).Related(&resultset)

	return resultset
}

//One Get one post
func (p *Post) One(db *gorm.DB, id int) Post {
	var result Post
	db.First(&result, "id = ?", id)

	return result
}

//Delete particular post
func (p *Post) Delete(db *gorm.DB) {
	db.Delete(p)
}
