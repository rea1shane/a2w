.PHONY: help run package build clean clean-build clean-package docker-run docker-build docker-rm-image

GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

APP_NAME = a2w
VERSION = $(shell cat VERSION)

PORT = 5001
TEMPLATE = ./templates/base.tmpl

BIN_DIR = bin
BIN_NAME = $(APP_NAME)
BIN_REFERENCE = $(BIN_DIR)/$(BIN_NAME)

PACKAGE_DIR = packages
PACKAGE_TEMP_DIR_NAME = $(APP_NAME)-$(VERSION).$(GOOS)-$(GOARCH)
PACKAGE_TEMP_DIR_REFERENCE = $(PACKAGE_DIR)/$(PACKAGE_TEMP_DIR_NAME)
PACKAGE_REFERENCE = $(PACKAGE_TEMP_DIR_REFERENCE).tar.gz

DOCKER_REPO = rea1shane
DOCKER_IMAGE_NAME = $(APP_NAME)

help: build
	$(BIN_REFERENCE) -h

run: build
	$(BIN_REFERENCE) --port $(PORT) --template $(TEMPLATE)

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
	docker run --name $(APP_NAME) -d -p $(PORT):5001 $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(VERSION) --template $(TEMPLATE)

docker-build:
	docker build -t $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(VERSION) .

docker-rm-image:
	docker rmi $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(VERSION)

helm-install: docker-build
	helm install $(APP_NAME) helm-chart

helm-uninstall:
	helm uninstall $(APP_NAME)
