package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gkatanacio/marvel-characters-api/internal/errs"
	"github.com/gkatanacio/marvel-characters-api/internal/marvel"
	"github.com/gorilla/mux"
)

type GetAllCharactersHandler struct {
	marvelService marvel.Servicer
}

func NewGetAllCharactersHandler(marvelService marvel.Servicer) *GetAllCharactersHandler {
	return &GetAllCharactersHandler{marvelService}
}

func (h *GetAllCharactersHandler) Handle(w http.ResponseWriter, r *http.Request) {
	charIds, err := h.marvelService.GetAllCharacterIds()
	if err != nil {
		log.Println(err)
		errorResponse(w, err)
		return
	}

	jsonResponse(w, charIds, http.StatusOK)
}

type GetCharacterInfoHandler struct {
	marvelService marvel.Servicer
}

func NewGetCharacterInfoHandler(marvelService marvel.Servicer) *GetCharacterInfoHandler {
	return &GetCharacterInfoHandler{marvelService}
}

func (h *GetCharacterInfoHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	charId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		errorResponse(w, errs.NewBadRequest("invalid id"))
		return
	}

	char, err := h.marvelService.GetCharacter(charId)
	if err != nil {
		log.Println(err)
		errorResponse(w, err)
		return
	}

	jsonResponse(w, char, http.StatusOK)
}
