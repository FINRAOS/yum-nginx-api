**[yum-nginx-api][1]**: A frontend API for yum repositories with NGINX
=======

yum-nginx-api is an easy frontend API for yum repositories running on the NGINX web server. It rapidly serves updates to Red Hat and CentOS and supports scaling.

It is a deployable solution with Docker or any existing web server with WSGI support. yum-nginx-api enables CI tools to be used for managing and promoting yum repositories.


**Problems solved with this project**:

1.  Serves updates to Red Hat / CentOS *really fast* and easily scalable.
2.  Limited options for a self-service yum repository to engineers via an API.
3.  Continuous Integration (CI) tools like Jenkins can build, sync, and promote yum repositories with this project unlike Red Hat Satellite Server and Spacewalk.
4.  Poor documentation on installing a yum repository with NGINX web server.



**Technologies Needed** [see install](#install) :

 1.  Server (Bare-metal/Cloud)
 2.  [NGINX Web Server][2]
 3.  [Python Flask][3] (Optional)
 4.  [Python Gunicorn][4] (Optional)
 5.  [Python Supervisor][5] (Optional)
 6.  [Docker][6] (Optional) 

**Diagram**:
![yum-nginx-api Diagram][7]


## Pull Docker image from the Docker Registry <a name="install"></a>
    docker pull finraos/yum-nginx-api
    docker run -d -p 80:80 finraos/yum-nginx-api

## How to build yum-nginx-api (Docker)

    git clone https://github.com/FINRAOS/yum-nginx-api.git
    cd yum-nginx-api && docker build -t finraos/yum-nginx-api .
    docker run -d -p 80:80 finraos/yum-nginx-api
    sleep 10 && docker logs `docker ps | grep 'yum-nginx-api' | awk '{ print $1 }' | head -n1` 

## How to Install yum-nginx-api (Non-Docker)

    # Need EPEL repo installed
    git clone https://github.com/FINRAOS/yum-nginx-api.git
    yum install -y python-pip supervisor gcc nginx createrepo python-setuptools
    cd yum-nginx-api && pip install -r requirements.txt
    mkdir -p /opt/repos/pre-release
    bash scripts/settings.sh
    cp -f supervisor/supervisord.conf /etc/
    cp -f supervisor/yumapi.conf /etc/supervisord.d/
    cp -rf yumapi /opt/
    cp -rf nginx/* /etc/nginx/
    supervisord -n -c /etc/supervisord.conf nohup &

## API Usage 

**Post binary RPM to API endpoint:**

    curl -F file=@yobot-4.6.2.noarch.rpm http://localhost/api/upload

**Successful post:**

    {
      "mime": "application/x-rpm", 
      "name": "yobot-4.6.2.noarch.rpm", 
      "size_mb": 294, 
      "status": 202
    }

**Unsuccessful post:**

    {
      "mime": "inode/x-empty", 
      "name": "yobot-4.6.2.noarch.rpm", 
      "size_mb": 0, 
      "status": 415
    }

  
## Contributing & Sponsor

More information on how to contribute to this project including sign off and the [DCO agreement](https://github.com/FINRAOS/yum-nginx-api/blob/master/DCO.md), please see the project's [GitHub wiki](https://github.com/FINRAOS/yum-nginx-api/wiki) for more information.

FINRA has graciously allocated time for their internal development resources to enhance yum-nginx-api and encourages participation in the open source community. Want to join FINRA? Please visit https://finra.org/careers.

[![FINRA Logo][8]](https://finra.org)


## License Type

yum-nginx-api project is licensed under [Apache License Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)


  [1]: http://github.com/finraos/yum-nginx-api/wiki
  [2]: http://nginx.org
  [3]: http://flask.pocoo.org
  [4]: http://gunicorn.org
  [5]: http://supervisord.org
  [6]: https://docker.io
  [7]: http://marshyski.com/yum-nginx-api.png
  [8]: http://www.finra.org/web/groups/corporate/@corp/documents/web_asset/p075334.gif
