BUILT_ON:=$(shell date)
COMMIT_HASH:=$(shell git log -n 1 --pretty=format:"%H")
PACKAGES:=$(shell go list ./... | sed -n '1!p' | grep -v /vendor/)
LDFLAGS:='-X "main.builtOn=$(BUILT_ON)" -X "main.commitHash=$(COMMIT_HASH)"'

default: run

test:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

cover: test
	go tool cover -html=coverage-all.out

run: config
	go run -ldflags $(LDFLAGS) *.go

build:
	go build -ldflags $(LDFLAGS) .

clean:
	rm -f yumapi yumapi.yaml coverage.out coverage-all.out

config:
	printf "upload_dir: .\ndev_mode: true" > yumapi.yaml

docker: build
	printf "upload_dir: /repo\n" > yumapi.yaml
	docker build -t finraos/yum-nginx-api .