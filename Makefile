.PHONY: run build clean docker-run docker-build docker-rm-image

APP_NAME = a2w
VERSION = $(shell cat VERSION)
BIN_REFERENCE = bin/$(APP_NAME)

run: build
	$(BIN_REFERENCE)

build: clean
	mkdir bin
	go build -v -o $(BIN_REFERENCE) ./...

clean:
	rm -rf bin

docker-run: docker-build
	docker run --name $(APP_NAME) -d -p 9099:9099 rea1shane/a2w:$(VERSION)

docker-build: docker-rm-image
	docker build -t rea1shane/a2w:$(VERSION) .

docker-rm-image:
	docker rmi rea1shane/a2w:$(VERSION)
