.PHONY: build test vet lint clean install run cross-build

BINARY_NAME=rs
MAIN_FILE=main.go

build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

test:
	go test -race -count=1 -coverprofile=coverage.out ./...

vet:
	go vet ./...

lint:
	golangci-lint run

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -rf dist/

install: build
	go install

run:
	go run $(MAIN_FILE)

cross-build:
	GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=arm64 go build -o build/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 go build -o build/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
