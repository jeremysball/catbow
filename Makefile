BINARY_NAME=cb

all: clean build test

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test ./catbow/

clean:
	go clean
	# does go clean -testcache do go clean? 
	go clean -testcache
	rm -f ${BINARY_NAME}
