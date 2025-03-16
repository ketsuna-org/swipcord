package main

import (
	"api/controllers"
	"api/models"
	"api/utils"
	"fmt"
	"log"
	"net/http"

	mux "github.com/gorilla/mux"
	"gorm.io/gorm"
)

var (
	database *gorm.DB
	initErr  error
)

func init() {
	// first we need to check if every environment variable is set
	// if not, we will panic and stop the execution
	envs := []string{
		"DB_HOST",
		"DB_USER",
		"DB_PASS",
		"DB_NAME",
		"DB_PORT",
		"DISCORD_CLIENT_ID",
		"DISCORD_CLIENT_SECRET",
	}

	if !utils.VerifyEnv(envs) {
		panic("Not all environment variables are set")
	}

	database, initErr = models.InitDatabase()
	if initErr != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", initErr))
	}

	// create Schema

	database.Exec("CREATE SCHEMA IF NOT EXISTS users")

	// Migrate the schema
	if database.AutoMigrate(&models.User{}, &models.Guilds{}, &models.History{}, &models.Reviews{}, &models.Statistics{}) != nil {
		panic("Failed to migrate the schema")
	} else {
		log.Println("Schema migrated successfully")
	}
}

func main() {
	// the Database successfully connected => do something now :)
	r := mux.NewRouter()
	// add Database to the context
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = models.ContextWithDatabase(ctx, database)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	r.Handle("/discord/login", http.HandlerFunc(controllers.DiscordOauth2))
	r.Handle("/discord/callback", http.HandlerFunc(controllers.DiscordCallback))

	log.Println("Server is running on port 4000")
	log.Panic(http.ListenAndServe(":4000", r))
}
