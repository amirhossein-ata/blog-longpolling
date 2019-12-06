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

//Author model
type Author struct {
	gorm.Model
	Name  string
	id    int    `gorm:"AUTO_INCREMENT"`
	Posts []Post `gorm:"foreignkey:AuthorID"`
}

//AuthorMethods - Author related methods
func AuthorMethods(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Content-Type", "aplication/json")
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
		// var author Author
		// _ = json.NewDecoder(r.Body).Decode(author)
		Body, _ := ioutil.ReadAll(r.Body)
		var authors []Author
		var author Author
		json.Unmarshal(Body, &author)
		fmt.Println(author.ID)
		db.Find(&authors, author.ID)
		if len(authors) != 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("User already created")

			return
		}

		db.Create(&author)

		json.NewEncoder(w).Encode(author)

	}

	if r.Method == http.MethodGet {
		var authors []Author
		vars := mux.Vars(r)
		user, err := strconv.Atoi(vars["user"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("id not entered corectly")
		}
		db.Find(&authors, user)
		if len(authors) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("User not found")
		}
		json.NewEncoder(w).Encode(authors)
	}

}
