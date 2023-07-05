package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main(){

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT is not found in the environment")
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

	router.Mount(
		"/v1", v1Router,
	)

	srv := &http.Server{
		Handler: router,
		Addr : ":"+portString,
	}
	err := srv.ListenAndServe()
	if err == nil{
		log.Fatal(err)
	}
}