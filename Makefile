.PHONY: run build docker-run docker-build clean

version="v0.1.0"

run: build
	bin/a2w

build: clean
	mkdir bin
	go build -v -o bin/a2w ./...

clean:
	rm -rf bin

docker-run:
	#TODO

docker-build:
	#TODO
