package main

import (
	"log"
	"net/http"

	"github.com/gkatanacio/marvel-characters-api/internal/handlers"
	"github.com/gkatanacio/marvel-characters-api/internal/marvel"
	"github.com/gorilla/mux"
)

// @title Marvel Characters API
// @version 1.0
// @description This API serves as a gateway for fetching character data from Marvel's API.
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
func main() {
	log.Println("server startup")

	cfg := marvel.NewConfig()
	client := marvel.NewClient(cfg)
	cache := marvel.NewInMemCache()
	service := marvel.NewService(client, cache)

	if cfg.EagerLoadCache {
		log.Println("prepopulating cache")
		if err := service.ReloadCache(); err != nil {
			log.Println(err)
			panic("failed to populate cache")
		}
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
