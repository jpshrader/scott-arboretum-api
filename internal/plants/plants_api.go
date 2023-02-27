package plants

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/jpshrader/scott-arboretum-api/response"
)

type Plant struct {
	CommonName     string `json:"commonName,omitempty"`
	ScientificName string `json:"scientificName,omitempty"`
}

func GetPlants(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("data/scott-arboretum-plant-list.json", os.O_RDONLY, fs.ModePerm)
	if err != nil {
		response.JsonEncode(w, http.StatusInternalServerError, "unable to open plant list")
		return
	}

	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		response.JsonEncode(w, http.StatusInternalServerError, "unable to read plant list")
		return
	}

	plants := []Plant{}
	err = json.Unmarshal(data, &plants)
	if err != nil {
		response.JsonEncode(w, http.StatusInternalServerError, "unable to parse plant list")
		return
	}

	response.JsonEncode(w, http.StatusOK, plants)
}
