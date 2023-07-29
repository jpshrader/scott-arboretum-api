package plants

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/jpshrader/scott-arboretum-api/response"
)

type Plant struct {
	CommonName     string `json:"commonName,omitempty"`
	ScientificName string `json:"scientificName,omitempty"`
	CommonType     string `json:"commonType,omitempty"`
	LookupName     string `json:"lookupName,omitempty"`
}

func GetPlants(w http.ResponseWriter, r *http.Request) {
	plants, err := readPlants(w)
	if err != nil {
		response.JsonEncode(w, http.StatusInternalServerError, err.Error())
		return
	}

	name := strings.ToLower(r.URL.Query().Get("name"))
	if name != "" {
		filteredPlants := []Plant{}
		for _, plant := range plants {
			if strings.Contains(strings.ToLower(plant.CommonName), name) {
				filteredPlants = append(filteredPlants, plant)
			}
		}
		plants = filteredPlants
	}

	response.JsonEncode(w, http.StatusOK, plants)
}

func readPlants(w http.ResponseWriter) ([]Plant, error) {
	file, err := os.OpenFile("data/scott-arboretum-plant-list.json", os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return []Plant{}, errors.New("unable to open plant list")
	}

	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return []Plant{}, errors.New("unable to read plant list")
	}

	plants := []Plant{}
	err = json.Unmarshal(data, &plants)
	if err != nil {
		return []Plant{}, errors.New("unable to parse plant list")
	}

	return plants, nil
}
