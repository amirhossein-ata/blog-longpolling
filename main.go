package main

import (
	"fmt"
	"log"
	"net/http"

	// "math/rand"
	// "strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func InitialMigration() {

	db, err := gorm.Open("sqlite3", "blog.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connnect to database")

	}
	defer db.Close()

	db.AutoMigrate(&Author{})

}

func main() {
	InitialMigration()
	r := mux.NewRouter()

	r.HandleFunc("/author/{user}", AuthorMethods).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	// r.HandleFunc("/post/{user}", Post).Methods("POST")
	r.Use(mux.CORSMethodMiddleware(r))
	log.Fatal(http.ListenAndServe(":8000", r))
}
