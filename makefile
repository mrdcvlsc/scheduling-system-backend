build_dev:
	go build -ldflags "-s -w" -o app && $(MAKE) frontend_local && ./app

build:
	go build -tags netgo -ldflags '-s -w' -o app

uml: project_uml.puml
	@echo UML already generated

project_uml.puml:
	@echo Generating UML...
	@echo using: https://github.com/jfeliu007/goplantuml
	rm project_uml.puml
	@goplantuml -recursive -show-aggregations -show-aliases -show-compositions -show-connection-labels -show-implementations -aggregate-private-members ./ > project_uml.puml

devr_local:
	rm -rf scheduling-system-temporary-data
	cd ../scheduling-system-temporary-data && python unify-subjects.py && python pack.py && mv release.zip ../scheduling-system-backend
	unzip release.zip -d ./scheduling-system-temporary-data

devr:
	@echo Download Test Data
	rm -rf scheduling-system-temporary-data
	gh release --repo github.com/mrdcvlsc/scheduling-system-temporary-data download --archive zip --clobber
	unzip scheduling-system-temporary-data-tmp-data-v*.zip -d ./ -x '*.py' '*.md' '*.git*' '*/makefile'
	mv scheduling-system-temporary-data-tmp-data-v*/ scheduling-system-temporary-data
	rm scheduling-system-temporary-data-tmp-data-v*.zip
	zip -r scheduling-system-temporary-data.zip scheduling-system-temporary-data

# old download temp data
# gh release --repo github.com/mrdcvlsc/scheduling-system-temporary-data download --pattern *.zip --clobber
# rm -rf scheduling-system-temporary-data
# mkdir scheduling-system-temporary-data
# unzip release.zip -d ./scheduling-system-temporary-data

frontend:
	@echo Download Frontend
	gh release --repo github.com/mrdcvlsc/scheduling-system-frontend download --pattern dist.zip --clobber
	rm -rf dist
	mkdir dist
	unzip dist.zip

frontend_local:
	rm -rf dist
	cd ../scheduling-system-frontend && npm run build && cp -R dist ../scheduling-system-backend

clean:
	go clean -testcache

test:
	go clean -testcache
	go test ./... -p 1 -timeout 0

testv:
	go clean -testcache
	go test ./... -v -p 1 -timeout 0

testvs:
	go clean -testcache && go test -run TestNewPopulation ./GeneticAlgorithm -v -timeout 0

bench:
	# we need to escape the dollar sign for the command: go test -run=^$ -bench=. ./...
	go test -run=^$$ -bench=. ./... -benchmem

todo:
	python todo.py

find:
	grep -nr "Persistence" ./

test_api:
	./app & sleep 1 && node Tests/schedule-serialization.js