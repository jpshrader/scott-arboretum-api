# starts the api server
run:
	go run main.go

# updates `data/scott-arboretum-plant-list.json` from the scott arboretum 'Plant Locator' website
list-plants:
	go run ./cmd/list_plants/main.go

build:
	go build -o bin/$(APP_NAME) main.go