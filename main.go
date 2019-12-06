package main

import (
	"fmt"
	"log"
	"net/http"

	// "math/rand"
	// "strconv"

	"github.com/gorilla/mux"
	"github.com/jcuga/golongpoll"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/cors"
)

func InitialMigration() {

	db, err := gorm.Open("sqlite3", "blog.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connnect to database")

	}
	defer db.Close()

	db.AutoMigrate(&Author{}, &Post{}, &Comment{})

}

func main() {
	manager, err := golongpoll.StartLongpoll(golongpoll.Options{
		LoggingEnabled: true,
	})
	if err != nil {
		log.Fatalf("Failed to create manager: %q", err)
	}
	go Longpoll(manager)
	InitialMigration()
	r := mux.NewRouter()
	r.HandleFunc("/post/{user}", Posts).Methods(http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPatch)
	r.HandleFunc("/author/{user}", AuthorMethods).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc("/comment/{post}", CommentMethods).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc("/like/{post}", Like).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc("/longpoll", manager.SubscriptionHandler)
	// r.Use(mux.CORSMethodMiddleware(r))
	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
