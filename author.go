package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//Author model
type Author struct {
	gorm.Model
	Name  string
	id    int
	Posts []Post `gorm:"foreignkey:AuthorID"`
}

//AuthorMethods - Author related methods
func AuthorMethods(w http.ResponseWriter, r *http.Request) {
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
		// var author Author
		// _ = json.NewDecoder(r.Body).Decode(author)
		Body, _ := ioutil.ReadAll(r.Body)
		var author Author
		json.Unmarshal(Body, &author)
		er := db.NewRecord(author)
		if er == false {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode("There is an author with that ID")

			return

		}
		db.Create(author)
		fmt.Println(author)

		json.NewEncoder(w).Encode(author)

	}

	if r.Method == http.MethodGet {
		var authors []Author
		var := Vars.
		db.Find(&authors, )
		json.NewEncoder(w).Encode(authors)
	}

}
