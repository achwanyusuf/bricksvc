Namespace = `echo bricksvc`
BuildTime = `date +%FT%T%z`
Version = `git describe --tag`

.PHONY: sqlboiler
sqlboiler:
	sqlboiler psql -c sqlboiler.toml --add-soft-deletes --add-global-variants

.PHONY: swaggo
swaggo: 
	go install github.com/swaggo/swag/cmd/swag@v1.16.3 && swag init -g ./src/cmd/main.go

.PHONY: build
build: swaggo ci run-tests
	go build -tags dynamic -ldflags "-X main.Namespace=${Namespace} -X main.BuildTime=${BuildTime}  -X main.Version=${Version}" -race -o ./build/app ./src/cmd

.PHONY: docker-build
docker-build: build
	sudo docker build -f script/Dockerfile -t ${Namespace}-${Version} .

.PHONY: docker-compose
docker-compose: build
	@sudo docker-compose down
	@sudo docker-compose pull
	@sudo docker-compose up --build -d

.PHONY: kill-process
kill-process:
	@lsof -i :8082 | awk '$$1 ~ /app/ { print $$2 }' | xargs kill -9 || true

.PHONY: run
run: kill-process build
	@./build/app

.PHONY: golangci-install
golangci-install:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.2
	@golangci-lint --version

.PHONY: ci
ci:
	$(shell go env GOPATH)/bin/golangci-lint run --verbose

.PHONY: mock
mock:
	@`go env GOPATH`/bin/mockgen -source src/repository/account/account.go -destination src/repository/mock/account/account.go
	@`go env GOPATH`/bin/mockgen -source src/repository/accountrole/accountrole.go -destination src/repository/mock/accountrole/accountrole.go
	@`go env GOPATH`/bin/mockgen -source src/repository/bank/bank.go -destination src/repository/mock/bank/bank.go
	@`go env GOPATH`/bin/mockgen -source src/repository/role/role.go -destination src/repository/mock/role/role.go
	@`go env GOPATH`/bin/mockgen -source src/repository/transfer/transfer.go -destination src/repository/mock/transfer/transfer.go
	@`go env GOPATH`/bin/mockgen -source src/usecase/account/account.go -destination src/usecase/mock/account/account.go
	@`go env GOPATH`/bin/mockgen -source src/usecase/accountrole/accountrole.go -destination src/usecase/mock/accountrole/accountrole.go
	@`go env GOPATH`/bin/mockgen -source src/usecase/bank/bank.go -destination src/usecase/mock/bank/bank.go
	@`go env GOPATH`/bin/mockgen -source src/usecase/role/role.go -destination src/usecase/mock/role/role.go
	@`go env GOPATH`/bin/mockgen -source src/usecase/transfer/transfer.go -destination src/usecase/mock/transfer/transfer.go

.PHONY: run-tests
run-tests:
	@GOEXPERIMENT=nocoverageredesign go test -v -tags dynamic `go list ./... | grep -i 'repository\|usecase'` -cover -failfast

.PHONY: generateproto
generateproto:
	@protoc --proto_path=conf/proto --go_out="./src/model/grpcmodel" --go_opt=paths=source_relative --go-grpc_out="./src/model/grpcmodel" --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional order.proto
