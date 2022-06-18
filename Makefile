
.DEFAULT_GOAL := build

BINARY=mdbrute.exe

BUILD=2022-06-18
TAGS=-tags forceposix
LDFLAGS=
#Stripping the exe causes the file to appear more suspicious, but can be re-enabled with:
#LDFLAGS=-ldflags "-w -w"

fetch-go:
	go mod download
	go mod tidy

upgrade-deps:
	go get -u ./...

upgrade-deps-patch:
	go get -u=patch ./...

build: fetch-go compile_default

manifest:
	rsrc -arch amd64 -ico resource/icon.ico -manifest mdbrute.exe.manifest

compile_default: fetch-go
	rm -f ${BINARY}
	env GOOS=windows CGO_ENABLED=0 GOARCH=amd64 go build ${LDFLAGS} ${TAGS} -o ${BINARY} .	
