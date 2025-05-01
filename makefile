SEARCH_PATTERN:=

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

gh_test:
	go clean -testcache
	go test ./... -p 1 -timeout 0 -skip "TestIntegration"

	@echo "running unit tests with fresh dev resources each"
	$(MAKE) devr
	go clean -testcache && go test -run Test ./GeneticAlgorithm -timeout 0
	$(MAKE) devr
	go clean -testcache && go test -run Test ./Resources/Curriculum -timeout 0
	$(MAKE) devr
	go clean -testcache && go test -run Test ./Resources/Instructors -timeout 0
	$(MAKE) devr
	go clean -testcache && go test -run Test ./Resources/Rooms -timeout 0
	$(MAKE) devr
	go clean -testcache && go test -run Test ./Schedule -timeout 0
	$(MAKE) devr
	go clean -testcache && go test -run Test ./StorageResources -timeout 0
	$(MAKE) devr
	go clean -testcache && go test -run Test ./StorageSchedule -timeout 0
	$(MAKE) devr
	go clean -testcache && go test -run Test ./Tests/schedule_datastructure_basic -timeout 0
	$(MAKE) devr
	go clean -testcache && go test -run Test ./Utils -timeout 0
	
	@echo "running integration tests"
	$(MAKE) devr
	go clean -testcache && go test -run TestIntegration ./ -timeout 0

gh_test_local:
	@echo "running unit tests with fresh dev resources each"
	$(MAKE) devr_local
	go clean -testcache && go test -run Test ./GeneticAlgorithm -timeout 0
	$(MAKE) devr_local
	go clean -testcache && go test -run Test ./Resources/Curriculum -timeout 0
	$(MAKE) devr_local
	go clean -testcache && go test -run Test ./Resources/Instructors -timeout 0
	$(MAKE) devr_local
	go clean -testcache && go test -run Test ./Resources/Rooms -timeout 0
	$(MAKE) devr_local
	go clean -testcache && go test -run Test ./Schedule -timeout 0
	$(MAKE) devr_local
	go clean -testcache && go test -run Test ./StorageResources -timeout 0
	$(MAKE) devr_local
	go clean -testcache && go test -run Test ./StorageSchedule -timeout 0
	$(MAKE) devr_local
	go clean -testcache && go test -run Test ./Tests/schedule_datastructure_basic -timeout 0
	$(MAKE) devr_local
	go clean -testcache && go test -run Test ./Utils -timeout 0

	@echo "running integration tests"
	$(MAKE) devr_local
	go clean -testcache && go test -run TestIntegration ./ -timeout 0

bench:
	# we need to escape the dollar sign for the command: go test -run=^$ -bench=. ./...
	go test -run=^$$ -bench=. ./... -benchmem

todo:
	python todo.py

find:
	grep -nr $(SEARCH_PATTERN) ./

test_api:
	./app & sleep 1 && node Tests/schedule-serialization.js