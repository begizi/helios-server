.PHONY: deps fmt build clean gin
export GOPATH=$(shell pwd)

deps:
	go get -d -v helios/...

fmt:
	go fmt helios/...

build: deps
	go install helios

clean:
	go clean -i -r helios/...

gin:
	cd src/helios && gin -i -a 8888 -p 8989
