build_dev:
	go build -ldflags "-s -w" -o app

build:
	go build -tags netgo -ldflags '-s -w' -o app

test:
	go test ./...