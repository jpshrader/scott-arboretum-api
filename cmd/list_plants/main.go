package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
)

const PLANT_LIST_URL = "http://arbnav.scottarboretum.org/weboi/oecgi3.exe/INET_ECM_GET_NAMELIST?time=1677530091407&NAMETYPE=COMM&BOXMODE=1"
const KEEP_DUPLICATES = false

type plantListItem struct {
	CommonName string `json:"commonName"`
}

func main() {
	log.Print("retrieving plant list...")

	res, err := http.Get(PLANT_LIST_URL)
	if err != nil {
		log.Fatal("unable to fetch plant list: ", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("unable to read plant list response: ", err)
	}

	plantListItems, err := getPlantListItems(string(body))
	if err != nil {
		log.Fatal("unable to parse plant list: ", err)
	}

	err = writePlantListItems(plantListItems)
	if err != nil {
		log.Fatal("unable to write plant list items: ", err)
	}

	log.Print("plant list retrieved")
}

func getPlantListItems(plantList string) ([]plantListItem, error) {
	plantNameList := map[string]bool{}
	plantListItems := []plantListItem{}

	regex := regexp.MustCompile(`<a href="javascript:void\(0\);".*?>(.{2,}?)<\/a>`)
	matches := regex.FindAllSubmatch([]byte(plantList), -1)

	for _, match := range matches {
		plantName := string(match[1])
		if _, found := plantNameList[plantName]; KEEP_DUPLICATES || !found {
			if KEEP_DUPLICATES && found {
				log.Printf("duplicate plant name: %s", plantName)
			}

			plantNameList[plantName] = true
			plantListItems = append(plantListItems, plantListItem{
				CommonName: plantName,
			})
		}
	}

	return plantListItems, nil
}

func writePlantListItems(plantListItems []plantListItem) error {
	file, err := json.MarshalIndent(plantListItems, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("data/scott-arboretum-plant-list.json", file, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
