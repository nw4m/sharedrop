APP     = sharedrop
VERSION = v1.0.0
LDFLAGS = -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: build run clean release

build:
	go build $(LDFLAGS) -o $(APP) ./

run:
	go run ./ $(ARGS)

clean:
	rm -rf dist/ $(APP) $(APP).exe

# build all platforms
release:
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(APP)-windows-amd64.exe ./
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/$(APP)-linux-amd64    ./
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/$(APP)-linux-arm64    ./
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/$(APP)-darwin-amd64   ./
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/$(APP)-darwin-arm64   ./
	@echo "Done! Binaries in dist/"
