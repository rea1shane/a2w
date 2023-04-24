.PHONY: run package build clean clean-build clean-package docker-run docker-build docker-rm-image

GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

APP_NAME = a2w
VERSION = $(shell cat VERSION)

BIN_DIR = bin
BIN_NAME = $(APP_NAME)
BIN_REFERENCE = $(BIN_DIR)/$(BIN_NAME)

PACKAGE_DIR = packages
PACKAGE_REFERENCE = $(PACKAGE_DIR)/$(APP_NAME)-$(VERSION).$(GOOS)-$(GOARCH).tar.gz

run: build
	$(BIN_REFERENCE)

package: build clean-package
	mkdir $(PACKAGE_DIR)
	tar zcvf $(PACKAGE_REFERENCE) templates -C $(BIN_DIR) $(BIN_NAME)

build: clean-build
	mkdir $(BIN_DIR)
	go build -v -o $(BIN_REFERENCE) ./...

clean: clean-build clean-package

clean-build:
	rm -rf $(BIN_DIR)

clean-package:
	rm -rf $(PACKAGE_DIR)

docker-run: docker-build
	docker run --name $(APP_NAME) -d -p 9099:9099 rea1shane/a2w:$(VERSION)

docker-build: docker-rm-image
	docker build -t rea1shane/a2w:$(VERSION) .

docker-rm-image:
	docker rmi rea1shane/a2w:$(VERSION)
