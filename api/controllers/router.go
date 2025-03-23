package controllers

import (
	"api/middlewares"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	r.Handle("/discord/login", http.HandlerFunc(DiscordOauth2))
	r.Handle("/discord/callback", http.HandlerFunc(DiscordCallback))
	r.Handle("/users/@me", middlewares.Chain(http.HandlerFunc(GetUser), middlewares.AuthMiddleware))

	log.Println("Server is running on port 4000")
	log.Panic(http.ListenAndServe(":4000", r))
}
