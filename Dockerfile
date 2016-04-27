FROM fedora:21
MAINTAINER Tim Marcinowski <marshyski@gmail.com>

USER root

RUN yum -y update
RUN yum install -y epel-release python-pip supervisor gcc nginx createrepo python-setuptools
RUN rm -rf /usr/share/nginx
ADD ./requirements.txt /tmp/requirements.txt
RUN pip install -r /tmp/requirements.txt
RUN yum remove -y gcc && yum clean all
RUN mkdir -p /opt/repos/pre-release
ADD nginx/nginx.conf /etc/nginx/nginx.conf
ADD nginx/mime.types /etc/nginx/mime.types
ADD yumapi /opt/yumapi
ADD supervisor/yumapi.conf /etc/supervisord.d/yumapi.conf
ADD supervisor/supervisord.conf /etc/supervisord.conf
ADD ./scripts/settings.sh /tmp/settings.sh
RUN /bin/bash /tmp/settings.sh

EXPOSE 80

CMD ["supervisord", "-n", "-c", "/etc/supervisord.conf"]
