list-plants:
	go run ./cmd/list_plants/main.go

run:
	go run main.go

build:
	go build -o bin/$(APP_NAME) main.go