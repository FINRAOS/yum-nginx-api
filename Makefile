PACKAGE_NAME:='yumapi'
BUILT_ON:=$(shell date)
COMMIT_HASH:=$(shell git log -n 1 --pretty=format:"%H")
PACKAGES:=$(shell go list ./... | grep -v /vendor/)
LDFLAGS:='-s -w -X "main.builtOn=$(BUILT_ON)" -X "main.commitHash=$(COMMIT_HASH)"'

default: docker

test:
	go test -cover -v $(PACKAGES)

update-deps:
	go get -u ./...
	go mod tidy

gofmt:
	go fmt ./...

lint: gofmt
	$(GOPATH)/bin/golint $(PACKAGES)
	$(GOPATH)/bin/golangci-lint run

run: config
	go run -ldflags $(LDFLAGS) `find . | grep -v 'test\|vendor\|repo' | grep \.go`

# Cross-compile from OS X to Linux using xgo
cc:
	xgo --targets=linux/amd64 -ldflags $(LDFLAGS) -out $(PACKAGE_NAME) .
	mv -f yumapi-* yumapi

# Build on Linux
build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -a -o $(PACKAGE_NAME) .

clean:
	rm -rf yumapi* coverage.out coverage-all.out repodata *.rpm *.sqlite

config:
	printf "upload_dir: .\ndev_mode: true" > yumapi.yaml

docker:
	printf "upload_dir: /repo\n" > yumapi.yaml
	docker build -t finraos/yum-nginx-api:latest .

# Run just API without NGINX
drun:
	docker run -d -p 8080:8080 --name yumapi finraos/yum-nginx-api

compose:
	docker-compose up