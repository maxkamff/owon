CURRENT_DIR=$(shell pwd)
APP=template
APP_CMD_DIR=./cmd

pull-submodule:
	git submodule update --init --recrusive

update-submodule:
	git submodule update --remote --merge

build:
	CGO_ENABLED=0 GOOS=darwin go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}

lint: ## Run golangci-lint with printing to stdout
	golangci-lint -c .golangci.yaml run --build-tags "musl" ./...

swag:
	echo ${REGISTRY}
	swag init -g api/router.go -o api/docs

run:
	go run cmd/main.go

migrate_up:
	migrate -path migrations/ -database postgres://maxkamff:12345@localhost:5432/user_db up

migrate_down:
	migrate -path migrations/ -database postgres://maxkamff:12345@localhost:5432/user_db down


migrate_force:
	migrate -path migrations/ -database postgres://maxkamff:12345@localhost:5432/user_db force 1


migrate_casbin:
	migrate create -ext sql -dir migrations -seq create_casbin_rule_table