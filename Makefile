ifndef OUTPUT
	OUTPUT=out/main
endif

init:
	go mod download all

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT} .