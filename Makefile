protogen:
	protoc -I ./proto/products/v1/ products.proto --go_out=plugins=grpc:./proto/products/v1/ products.proto

run: protogen
	docker-compose -f docker-compose.dev.yml up --force-recreate --remove-orphans -d
	set -a && . ./.env.dev && set +a && go run .

build: protogen
	go build .

SCALE ?= 2

deploy: protogen
	docker-compose up --force-recreate --remove-orphans -d --scale grpc-atlant='$(SCALE)'

