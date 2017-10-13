PACKAGE_NAME:='yumapi'
BUILT_ON:=$(shell date)
COMMIT_HASH:=$(shell git log -n 1 --pretty=format:"%H")
PACKAGES:=$(shell go list ./... | sed -n '1!p' | grep -v /vendor/)
LDFLAGS:='-X "main.builtOn=$(BUILT_ON)" -X "main.commitHash=$(COMMIT_HASH)"'

default: docker

test:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

cover: test
	go tool cover -html=coverage-all.out

run: config
	go run -ldflags $(LDFLAGS) *.go

# Cross-compile from OS X to Linux using xgo
cc:
	xgo --targets=linux/amd64 -ldflags $(LDFLAGS) -out $(PACKAGE_NAME) .
	mv -f yumapi-* yumapi

# Build on Linux
build:
	CGO_ENABLED="1" go build -ldflags $(LDFLAGS) -o $(PACKAGE_NAME) .

clean:
	rm -rf yumapi* coverage.out coverage-all.out repodata *.rpm *.sqlite

config:
	printf "upload_dir: .\ndev_mode: true" > yumapi.yaml

docker: cc
	printf "upload_dir: /repo\n" > yumapi.yaml
	docker build -t finraos/yum-nginx-api:latest .

# Run just API without NGINX
drun:
	docker run -d -p 8080:8080 --name yumapi finraos/yum-nginx-api

compose:
	docker-compose up