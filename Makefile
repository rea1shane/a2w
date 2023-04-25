.PHONY: run package build clean clean-build clean-package docker-run docker-build docker-rm-image

GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

APP_NAME = a2w
VERSION = $(shell cat VERSION)

ARGS = --port $(PORT) --template $(TEMPLATE)
PORT = 5001
TEMPLATE = ./templates/base.tmpl

BIN_DIR = bin
BIN_NAME = $(APP_NAME)
BIN_REFERENCE = $(BIN_DIR)/$(BIN_NAME)

PACKAGE_DIR = packages
PACKAGE_TEMP_DIR_NAME = $(APP_NAME)-$(VERSION).$(GOOS)-$(GOARCH)
PACKAGE_TEMP_DIR_REFERENCE = $(PACKAGE_DIR)/$(PACKAGE_TEMP_DIR_NAME)
PACKAGE_REFERENCE = $(PACKAGE_TEMP_DIR_REFERENCE).tar.gz

run: build
	$(BIN_REFERENCE) $(ARGS)

package: build clean-package
	mkdir -p $(PACKAGE_TEMP_DIR_REFERENCE)
	cp $(BIN_REFERENCE) $(PACKAGE_TEMP_DIR_REFERENCE)
	cp -R templates $(PACKAGE_TEMP_DIR_REFERENCE)
	tar zcvf $(PACKAGE_REFERENCE) -C $(PACKAGE_DIR) $(PACKAGE_TEMP_DIR_NAME)
	rm -rf $(PACKAGE_TEMP_DIR_REFERENCE)

build: clean-build
	mkdir $(BIN_DIR)
	go build -v -o $(BIN_REFERENCE) ./...

clean: clean-build clean-package

clean-build:
	rm -rf $(BIN_DIR)

clean-package:
	rm -rf $(PACKAGE_DIR)

docker-run: docker-build
	docker run --name $(APP_NAME) -d -p 5001:5001 rea1shane/a2w:$(VERSION) $(ARGS)

docker-build: docker-rm-image
	docker build -t rea1shane/a2w:$(VERSION) .

docker-rm-image:
	docker rmi rea1shane/a2w:$(VERSION)
