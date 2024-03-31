package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jpshrader/scott-arboretum-api/internal/plants"
	"github.com/jpshrader/scott-arboretum-api/response"
)

const baseUrl = "/api"

func main() {
	log.Println("starting...")
	http.HandleFunc(fmt.Sprintf("GET %s", baseUrl), root)

	http.HandleFunc(fmt.Sprintf("GET %s/plants", baseUrl), plants.GetPlants)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func root(w http.ResponseWriter, r *http.Request) {
	response.Ok(w, "welcome to the unofficial scott arboretum api")
}
