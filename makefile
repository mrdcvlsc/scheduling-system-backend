build_dev:
	go build -ldflags "-s -w" -o app

build:
	go build -tags netgo -ldflags '-s -w' -o app

devr:
	gh release --repo github.com/mrdcvlsc/scheduling-system-temporary-data download --pattern *.zip --clobber
	mkdir scheduling-system-temporary-data
	unzip release.zip -d ./scheduling-system-temporary-data

test:
	go clean -testcache
	go test ./...

testv:
	go clean -testcache
	go test ./... -v

testvs:
	go clean -testcache
	go test -run TestJsonFilePersistence_GetAllRoom github.com/mrdcvlsc/scheduling-system-backend/Storage -v

bench:
	# we need to escape the dollar sign for the command: go test -run=^$ -bench=. ./...
	go test -run=^$$ -bench=. ./... -benchmem