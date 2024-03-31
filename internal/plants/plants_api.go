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

type plant struct {
	Name               string   `json:"name"`
	CommonName         string   `json:"commonName"`
	SortName           string   `json:"sortName"`
	PlantCateory       string   `json:"plantCategory"`
	ArboretumLocations []string `json:"arboretumLocations"`
}

func GetPlants(w http.ResponseWriter, r *http.Request) {
	plants, err := readPlants()
	if err != nil {
		response.InteralServerError(w, err.Error())
		return
	}

	commonName := strings.ToLower(r.URL.Query().Get("commonName"))
	if commonName != "" {
		filteredPlants := []plant{}
		for _, plant := range plants {
			if strings.Contains(strings.ToLower(plant.CommonName), commonName) {
				filteredPlants = append(filteredPlants, plant)
			}
		}
		plants = filteredPlants
	}

	response.Ok(w, plants)
}

func readPlants() ([]plant, error) {
	file, err := os.OpenFile("data/scott-arboretum-plants.json", os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return []plant{}, errors.New("unable to open plant list")
	}

	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return []plant{}, errors.New("unable to read plant list")
	}

	plants := []plant{}
	err = json.Unmarshal(data, &plants)
	if err != nil {
		return []plant{}, errors.New("unable to decode plant list")
	}

	return plants, nil
}
