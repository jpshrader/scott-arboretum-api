package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jpshrader/scott-arboretum-api/internal/plants"
	"github.com/jpshrader/scott-arboretum-api/response"
)

const BASE_URL = "/api"

func main() {
    log.Println("starting...")
    http.HandleFunc(BASE_URL, root)

    http.HandleFunc(fmt.Sprintf("%s/plants", BASE_URL), plants.GetPlants)


    log.Fatal(http.ListenAndServe(":8080", nil))
}

func root(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    response.JsonEncode(w, http.StatusOK, "Welcome to the Scott Arboretum Api!")
}