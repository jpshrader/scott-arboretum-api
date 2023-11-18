# starts the api server
run:
	go run main.go

# updates `data/scott-arboretum-plants.json` from the scott arboretum website
fetch-plants:
	go run ./cmd/fetch_plants/main.go

build:
	go build -o bin/$(APP_NAME) main.go

test:
	go test -v ./...

update:
	go get -u ./...
	go mod tidy