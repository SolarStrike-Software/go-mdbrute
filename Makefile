
.DEFAULT_GOAL := build

BINARY=mdbrute.exe

fetch-go:
	go mod download
	go mod tidy

upgrade-deps:
	go get -u ./...

upgrade-deps-patch:
	go get -u=patch ./...

build: fetch-go compile_default

compile_default: fetch-go
	rm -f ${BINARY}
	env GOOS=windows CGO_ENABLED=0 GOARCH=amd64 go build -ldflags '-s -w' -tags forceposix -o ${BINARY} .	
