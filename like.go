package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func Like(w http.ResponseWriter, r *http.Request) {
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
	if r.Method == http.MethodPost {
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
		var p Post
		db.Where("id = ?", post).Find(&p)
		p.Likes++
		db.Save(&p)
		db.Where("id = ?", post).Find(&p)
		json.NewEncoder(w).Encode(p)
	}
}
