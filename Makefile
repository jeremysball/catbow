BINARY_NAME=catbow
BIN_DIR=bin

all: clean build test

build: linux windows
linux:
	GOOS=linux go build -o ${BIN_DIR}/${BINARY_NAME} main.go
windows:
	GOOS=windows go build -o ${BIN_DIR}/${BINARY_NAME}.exe main.go

test:
	go test ./catbow/

clean:
	go clean
	# does go clean -testcache do go clean?
	go clean -testcache
	rm -rf ${BIN_DIR}
