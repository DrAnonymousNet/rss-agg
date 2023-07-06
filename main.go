package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/dranonymousnet/rss-agg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}


func main(){

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT is not found in the environment")
	}


	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
		log.Fatal("Unable to retrieve Database URL in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database")
	}
	apiCfg := apiConfig{
		DB: database.New(conn),
	} 

	log.Printf("Server starting on %v", portString)
	router := chi.NewRouter()

	router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "DELETE","PUT", "PATCH", "OPTIONS"},
			AllowedHeaders: []string{"*"},
			ExposedHeaders: []string{"link"},
			AllowCredentials: false,
			MaxAge: 300,
		}),
	)

	v1Router := chi.NewRouter()

	v1Router.Get(
		"/healthz", handlerReadiness,
	)
	v1Router.Get("/err", handlerError,
	)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount(
		"/v1", v1Router,
	)

	srv := &http.Server{
		Handler: router,
		Addr : ":"+portString,
	}
	err = srv.ListenAndServe()
	if err == nil{
		log.Fatal(err)
	}
}