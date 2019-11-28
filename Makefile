.PHONY: build clean tool lint help
all: build

build:
    go build -v .

tool:
    go tool vet . |& grep -v vendor;true
    gofmt -w .build

lint:
    golint ./..

clean:
    rm -rf gin-blog
    go clean -i .

help:
    @echo "make :compile packages add dependencies"
    @echo "make tool: run specified go tool"
    @echo "make lint: golint ./.."
    @echo "make clean: remove object files and cache files"
