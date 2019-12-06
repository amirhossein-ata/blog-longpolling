package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//Post Model
type Post struct {
	gorm.Model
	PostID   uint `gorm:"AUTO_INCREMENT"`
	Title    string
	Text     string
	Likes    uint
	AuthorID uint
	Comments []Comment `gorm:"foreignkey:PostID"`
}

//Posts - post related operations
func Posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "aplication/json")
	db, err := gorm.Open("sqlite3", "blog.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connnect to database")

	}
	defer db.Close()

	if r.Method == http.MethodOptions {
		return
	}
	//Get posts of a user
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		var authors []Author
		user, err := strconv.Atoi(vars["user"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("id not entered corectly")

		}
		db.Find(&authors, user)
		if len(authors) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("User not found")

			return
		}
		var posts []Post
		db.Where("author_id = ?", user).Find(&posts)

		json.NewEncoder(w).Encode(posts)

		return
	}
	if r.Method == http.MethodPost {
		Body, _ := ioutil.ReadAll(r.Body)
		var post Post
		json.Unmarshal(Body, &post)

		db.Create(post)
		json.NewEncoder(w).Encode(post)

	}
	if r.Method == http.MethodPatch {

		vars := mux.Vars(r)
		var posts []Post
		post, err := strconv.Atoi(vars["user"])
		fmt.Println(post)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("id not entered corectly")

			return

		}
		db.Find(&posts, post)
		if len(posts) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("post not found")

			return
		}
		var p Post
		db.Where("id = ?", post).Find(&p)

		Body, _ := ioutil.ReadAll(r.Body)
		var p1 Post
		json.Unmarshal(Body, &p1)

		p.Text = p1.Text
		db.Save(&p)
		db.Where("id = ?", post).Find(&p)
		json.NewEncoder(w).Encode(p)
	}

}
