<img src="/img/logo.png" alt="drawing" width="520px"/>

A GO API for managing yum repos and  NGINX to serve them
=======

[![CircleCI](https://circleci.com/gh/FINRAOS/yum-nginx-api/tree/master.svg?style=svg)](https://circleci.com/gh/FINRAOS/yum-nginx-api/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/FINRAOS/yum-nginx-api)](https://goreportcard.com/report/github.com/FINRAOS/yum-nginx-api)

[yum-nginx-api][1] is a go API for uploading RPMs to yum repositories and also configurations for running NGINX to serve them.

It is a deployable solution with Docker or a single 8MB statically linked Linux binary. yum-nginx-api enables CI tools to be used for uploading RPMs and managing yum repositories.

Included in this project is a go package `repojson` that can be used to read a repodata directory and return a JSON array of all packages in the primary.sqlite.(bz2|xz).  For usage go to [RepoJSON](#repojson)

**Problems solved with this project**:

1.  Serves updates to Red Hat / CentOS *really fast* and easily scalable.
2.  Limited options for a self-service yum repository to engineers via an API.
3.  Continuous Integration (CI) tools like Jenkins can build, sync, and promote yum repositories with this project unlike Red Hat Satellite Server and Spacewalk.
4.  Poor documentation on installing a yum repository with NGINX.

**Requirements**:

 1.  Server (Bare-metal/VM)
 2.  [NGINX][2]
 3.  [Go][3] >= 1.9.1 (Optional)
 4.  [xgo][4] (Optional)
 5.  [Docker][5] >=17.09/1.32 (Optional) 
 6.  [Docker Compose][6] >=1.16.1 (Optional)


## Run only API

    docker run -d -p 8080:8080 --name yumapi finraos/yum-nginx-api

## Run Docker Compose
	
	docker-compose up

## How to build yum-nginx-api (Go & Docker)

    # This projects needs CGO if not on Linux
    git clone https://github.com/FINRAOS/yum-nginx-api.git
    cd yum-nginx-api
    make build # Linux only
    make cc # OS X only
    make docker

## How to Install yum-nginx-api (Binary)

    make build # Linux only
    make cc # OS X only

## Configuration File `yumapi.yaml`

**Configuration file can be JSON, TOML, YAML, HCL, or Java properties**

    # createrepo workers, default is 2
    createrepo_workers:
    # http max content upload, default is 10000000 <- 10MB
    max_content_length:
    # yum repo directory, default is ./
    upload_dir:
    # port to run http server, default is 8080
    port:
    # max retries to retry failed createrepo, default is 3
    max_retries:

## API Usage 

**Post binary RPM to API endpoint:**

    curl -F file=@yobot-4.6.2.noarch.rpm http://localhost:8080/api/upload

**List repo contents package name, arch, version and summary:**

    curl http://localhost:8080/api/repo

**Successful post:**

    [{
        "name": "yobot",
        "arch": "x86_64",
        "version": "4.6.2",
        "summary": "Shenandoah RPM"
    },{
        "name": "yum-nginx-api-test",
        "arch": "x86_64",
        "version": "0.1",
        "summary": "Yum NGINX API Test RPM"
    }]
 
**Health check API endpoint**
 
    curl http://localhost/api/health

## RepoJSON

    package main

    import (
	    "encoding/json"
	    "fmt"

	    "github.com/FINRAOS/yum-nginx-api/repojson"
    )

    func main() {
	    ar, err := repojson.RepoJSON("./")
	    if err != nil {
		    fmt.Println(err)
	    }
	    js, err := json.Marshal(ar)
	    if err != nil {
		    fmt.Println(err)
	    }
	    fmt.Println(string(js))
    }

## Contributing & Sponsor

More information on how to contribute to this project including sign off and the [DCO agreement](https://github.com/FINRAOS/yum-nginx-api/blob/master/DCO.md), please see the project's [GitHub wiki](https://github.com/FINRAOS/yum-nginx-api/wiki) for more information.

FINRA has graciously allocated time for their internal development resources to enhance yum-nginx-api and encourages participation in the open source community. Want to join FINRA? Please visit [http://technology.finra.org/careers.html](http://technology.finra.org/careers.html).

[FINRA Technology](http://technology.finra.org/)

## License Type

yum-nginx-api project is licensed under [Apache License Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

  [1]: https://github.com/finraos/yum-nginx-api/wiki
  [2]: https://nginx.org
  [3]: https://golang.org
  [4]: https://github.com/karalabe/xgo
  [5]: https://docs.docker.com/engine/installation/
  [6]: https://docs.docker.com/compose/install/
