.PHONY: docker clean

build: $(shell find . -iname '*.go')
	go build -o bin/go-webserver main.go
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o bin/go-webserver-dist main.go

docker: Dockerfile bin/go-webserver-dist
	docker image build -t rhysemmas/go-webserver:latest .

clean:
	rm -rf bin
