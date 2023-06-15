CURRENT_DIR=$(shell pwd)

pull-submodule:
	git submodule update --init --recrusive

update-submodule:
	git submodule update --remote --merge

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}
	ls genproto/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp,}"

run:
	go run cmd/main.go

migrate_up:
	migrate -path migrations/ -database postgres://maxkamff:12345@localhost:5432/post_db up

migrate_down:
	migrate -path migrations/ -database postgres://maxkamff:12345@localhost:5432/post_db down


migrate_force:
	migrate -path migrations/ -database postgres://maxkamff:12345@localhost:5432/post_db force 2


migrate_posts:
	migrate create -ext sql -dir migrations -seq create_posts_table		