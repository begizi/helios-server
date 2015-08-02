.PHONY: deps fmt build clean serve
export GOPATH=$(shell pwd)

deps:
	go get -d -v helios/...

fmt:
	go fmt helios/...

build: deps
	go install helios

clean:
	go clean -i -r helios/...

serve:
ifeq ("$(wildcard $(./bin/fresh))","")
	@go get github.com/pilu/fresh
	@go install github.com/pilu/fresh
endif
	@./bin/fresh -c fresh.conf
