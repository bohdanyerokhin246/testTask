// main.go
package main

import (
	"net/http"
	"testTask/handlers"
	"testTask/middleware"
	"testTask/models"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db := models.InitDB()
	handlers.Init(db)

	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	r.Handle("/images", middleware.AuthMiddleware(http.HandlerFunc(handlers.UploadImage))).Methods("POST")
	r.Handle("/images", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetImages))).Methods("GET")

	handler := cors.Default().Handler(r)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		return
	}
}
