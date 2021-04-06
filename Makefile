.PHONY: clean build vendor iterate
export COMPOSE_FILE=docker-compose.yaml:docker-compose.build.yaml

all: clean vendor build
iterate: build restart

clean:
	rm -rf walker.so vendor

vendor:
	docker-compose run --rm go \
		mod vendor

build:
	docker-compose run --rm go \
		build -buildmode=plugin -trimpath -o ./dst/walker.so .

restart:
	docker-compose restart nakama