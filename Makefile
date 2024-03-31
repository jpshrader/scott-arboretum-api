run:
	go run main.go

fetch-plants:
	go run ./cmd/fetch_plants/main.go \
		--githubToken=$(githubToken)

build:
	go build -v ./...

test:
	go test -v ./...

update:
	go get -u ./...
	go mod tidy