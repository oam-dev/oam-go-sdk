
GO111MODULE := on

all: clean build

clean:
	rm -rf ./bin

build:
	CGO_ENABLED=0 \
	GO111MODULE=$(GO111MODULE) \
	go build -i -v -o ./bin/oambuilder main.go

