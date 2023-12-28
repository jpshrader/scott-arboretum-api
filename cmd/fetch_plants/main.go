package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/google/go-github/v57/github"
)

const (
	owner             = "jpshrader"
	slug              = "scott-arboretum-api"
	plantListPath     = "./data/scott-arboretum-plants.json"
	arboretumPlantUrl = "https://silva.swarthmore.edu/server/rest/services/Plant_Centers_Public_View/MapServer/1/query?f=json&resultOffset=0&resultRecordCount=9999&where=NAME%20IS%20NOT%20NULL&orderByFields=NAME&outFields=OBJECTID%2CACC_NUM_AND_QUAL%2CNAME%2CCOMMON_NAME_PRIMARY%2CDESCRIPTOR%2CCV_GROUP%2CETI%2CSORT_NAME%2CHABIT_FULL%2CSPEC_CHAR_FULL%2CCURRENT_LOCATION_FULL%2CCOMMON_NAME_PRIMARY&returnGeometry=false&spatialRel=esriSpatialRelIntersects"
)

type arboretumPlantPayload struct {
	Features []struct {
		Attributes struct {
			Name              string `json:"NAME"`
			CommonName        string `json:"COMMON_NAME_PRIMARY"`
			SortName          string `json:"SORT_NAME"`
			PlantCateory      string `json:"HABIT_FULL"`
			ArboretumLocation string `json:"CURRENT_LOCATION_FULL"`
		} `json:"attributes"`
	} `json:"features"`
}

type plant struct {
	Name               string   `json:"name"`
	CommonName         string   `json:"commonName"`
	SortName           string   `json:"sortName"`
	PlantCateory       string   `json:"plantCategory"`
	ArboretumLocations []string `json:"arboretumLocations"`
}

func main() {
	var githubToken string
	flag.StringVar(&githubToken, "githubToken", "", "api token with read/write access to this repo")
	flag.Parse()
	log.Print("retrieving plants...")

	res, err := http.Get(arboretumPlantUrl)
	if err != nil {
		log.Fatal("unable to fetch plant list: ", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("unable to read plant list response: ", err)
	}

	var plantPayload arboretumPlantPayload
	err = json.Unmarshal(body, &plantPayload)
	if err != nil {
		log.Fatal("unable to unmarshal plant list response: ", err)
	}

	plants, err := processPayload(plantPayload)
	if err != nil {
		log.Fatal("unable to parse plant list: ", err)
	}

	log.Println("found", len(plants), "plants")
	slices.SortFunc(plants, func(i, j plant) int {
		if i.SortName < j.SortName {
			return -1
		} else if i.SortName > j.SortName {
			return 1
		}
		return 0
	})

	err = writePlantListItems(plants)
	if err != nil {
		log.Fatal("unable to write plant list items: ", err)
	}

	if githubToken != "" {
		err = updatePlantList(context.Background(), githubToken, plants)
		if err != nil {
			log.Fatal("unable to update plant list: ", err)
		}
	}

	log.Print("plants retrieved")
}

func processPayload(payload arboretumPlantPayload) ([]plant, error) {
	plantLookup := make(map[string]plant, len(payload.Features))
	for _, f := range payload.Features {
		p, found := plantLookup[f.Attributes.Name]
		if found {
			if slices.Contains(p.ArboretumLocations, f.Attributes.ArboretumLocation) {
				continue
			}
			p.ArboretumLocations = append(p.ArboretumLocations, f.Attributes.ArboretumLocation)
			plantLookup[f.Attributes.Name] = p
		} else {
			plantLookup[f.Attributes.Name] = plant{
				Name:               f.Attributes.Name,
				CommonName:         f.Attributes.CommonName,
				SortName:           f.Attributes.SortName,
				PlantCateory:       f.Attributes.PlantCateory,
				ArboretumLocations: []string{f.Attributes.ArboretumLocation},
			}
		}
	}

	plants := make([]plant, 0, len(plantLookup))
	for _, plant := range plantLookup {
		plants = append(plants, plant)
	}

	return plants, nil
}

func writePlantListItems(plants []plant) error {
	file, err := json.MarshalIndent(plants, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("data/scott-arboretum-plants.json", file, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func updatePlantList(ctx context.Context, token string, plants []plant) error {
	client := github.NewClient(nil).WithAuthToken(token)

	plantListSha, err := getPlanListSha(ctx, client)
	if err != nil {
		return err
	}

	updatedPlantList, err := json.MarshalIndent(plants, "", "    ")
	if err != nil {
		return err
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("updated arboretum plant list"),
		Content: updatedPlantList,
		SHA:     &plantListSha,
	}
	_, _, err = client.Repositories.UpdateFile(ctx, owner, slug, plantListPath, opts)
	if err != nil {
		return err
	}
	return nil
}

func getPlanListSha(ctx context.Context, client *github.Client) (string, error) {
	file, _, _, err := client.Repositories.GetContents(ctx, owner, slug, plantListPath, nil)
	if err != nil {
		return "", err
	}
	return file.GetSHA(), nil
}
