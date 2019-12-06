package main

import (
	"fmt"
	"log"
	"net/http"

	// "math/rand"
	// "strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jcuga/golongpoll"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/cors"
)

func addCorsHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		fn(w, r)
	}
}

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
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r := mux.NewRouter()
	r.HandleFunc("/post/{user}", Posts).Methods(http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPatch)
	r.HandleFunc("/author/{user}", AuthorMethods).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc("/comment/{post}", CommentMethods).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc("/like/{post}", Like).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc("/longpoll", addCorsHeaders(manager.SubscriptionHandler))
	// r.Use(mux.CORSMethodMiddleware(r))
	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(origins, headers, methods)(handler)))
}
