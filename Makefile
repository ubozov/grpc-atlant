run:
	protoc -I ./proto/products/v1/ products.proto --go_out=plugins=grpc:./proto/products/v1/ products.proto
	go run .

build:
	protoc -I ./proto/products/v1/ products.proto --go_out=plugins=grpc:./proto/products/v1/ products.proto
	go build .

SCALE ?= 2

deploy:
	docker-compose up --force-recreate --remove-orphans -d --scale grpc-atlant='$(SCALE)'

