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

//Comment model
type Comment struct {
	Text   string
	PostID uint
}

//CommentMethods - Comment related methods
func CommentMethods(w http.ResponseWriter, r *http.Request) {

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
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		var posts []Post
		post, err := strconv.Atoi(vars["post"])

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("id not entered corectly")

		}
		db.Find(&posts, post)
		if len(posts) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("post not found")

			return
		}
		var comments []Comment

		db.Where("post_id = ?", post).Find(&comments)
		json.NewEncoder(w).Encode(comments)

		return

	}
	if r.Method == http.MethodPost {
		Body, _ := ioutil.ReadAll(r.Body)
		var comment Comment
		json.Unmarshal(Body, &comment)

		db.Create(&comment)
		json.NewEncoder(w).Encode(comment)

	}

}
