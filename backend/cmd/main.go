package main

import (
	"backend/internal/algorithm"
	"backend/internal/api"
	db "backend/internal/database"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func initAPIHandler() (api.Handler, error) {
	log.Print("Started!")

	database := db.NewDatabaseConnection("redis:6379", "")
	err := database.Connect()
	if err != nil {
		return api.Handler{}, err
	}
	log.Print("connected to db")

	alg1 := algorithm.NewAlgorithm("alg1", "http://algorithm:80")
	algorithms := []api.IAlgorithm{alg1}
	apiHandler := api.NewHandler(database, algorithms)
	if err := apiHandler.Config(); err != nil {
		return api.Handler{}, err
	}
	return apiHandler, nil
}

func main() {
	router := mux.NewRouter()

	apiHandler, err := initAPIHandler()
	if err != nil {
		log.Print("couldn't connect to database")
		return
	}
	apiHandler.EndpointInitialize(router)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "PUT", "GET"},
	})

	handler := c.Handler(router)
	port := "8081"
	log.Printf("Listening on port %s!", port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Failed to listen")
	}
}
