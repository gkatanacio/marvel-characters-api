package main

import (
	"log"
	"net/http"

	"github.com/gkatanacio/marvel-characters-api/internal/handlers"
	"github.com/gkatanacio/marvel-characters-api/internal/marvel"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("server startup")

	cfg := marvel.NewConfig()
	client := marvel.NewClient(cfg)
	cache := marvel.NewInMemCache()
	service := marvel.NewService(client, cache)

	log.Println("prepopulating cache")
	if err := service.ReloadCache(); err != nil {
		log.Println(err)
		panic("failed to populate cache")
	}

	getAllCharactersHandler := handlers.NewGetAllCharactersHandler(service)
	getCharacterInfoHandler := handlers.NewGetCharacterInfoHandler(service)

	r := mux.NewRouter()
	r.HandleFunc("/characters", getAllCharactersHandler.Handle).Methods(http.MethodGet)
	r.HandleFunc("/characters/{id}", getCharacterInfoHandler.Handle).Methods(http.MethodGet)

	port := ":8080"
	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
