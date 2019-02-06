
.PHONY:build
build:
	# build app for a linux based container
	@CGO_ENABLED=0 GOOS=linux go build -o ./payment -a -ldflags '-s' -installsuffix cgo ./cmd/payment/main.go

.PHONY:install
install:
	# install dependencies from gopkg file
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/go-swagger/go-swagger/cmd/swagger
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

.PHONY: test
test:
	# launch test across all project
	@go test -race ./...
	@go list ./... | grep -v /vendor/ | grep -v pb | xargs -L1 golint -set_exit_status
	@go vet `go list ./... | grep -v /vendor/`

.PHONY:coverage
coverage:
	# use go ability to generate an html with test coverage
	@go test `go list ./... | grep -v /vendor/` -cover -coverprofile=cover.out
	@go tool cover -html=cover.out

.PHONY:dev
dev:
	# launch external dependencies from docker-compose for local development
	@docker-compose up db adminer jaeger 

.PHONY:build
local: build
	# launch complete docker-compose for local execution
	@docker-compose build
	@docker-compose up

.PHONY:stop
stop:
	# shutdown docker-compose and clean payment image
	@docker-compose down
	@docker rmi payment

.PHONY:run
run:
	# launch payment-api 
	@go run cmd/payment/main.go

.PHONY:migration
migration:
	# launch migration.go to create/update schema and insert mocked data
	@go run cmd/migration/migration.go 

# .PHONY:doc
# doc: build
# 	# generate api documentation using go-swagger
# 	@swagger generate spec -o ./docs/swaggerui/swagger.json -b ./cmd/payment/ --scan-models
# 	@swagger validate ./docs/swaggerui/swagger.json

.PHONY: update-mocks
update-mocks:
	@go get github.com/vektra/mockery/.../
	@go list -f '{{.Dir}}' ./... | grep -v "payment-api$$" | xargs -n1 ${GOPATH}/bin/mockery -inpkg -case "underscore" -all -note "NOTE: run 'make update-mocks' from payment-api top folder to update this file and generate new ones." -dir || true

.PHONY: integration
integration:
	# launch test flagged as integration
	@CGO_ENABLED=0 go test -v `go list ./... | grep -v /vendor/` -tags=integration