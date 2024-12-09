build_dev:
	go build -ldflags "-s -w" -o app

build:
	go build -tags netgo -ldflags '-s -w' -o app

test:
	go clean -testcache
	go test ./Tests/...

testv:
	go clean -testcache
	go test ./Tests/... -v