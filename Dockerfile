FROM centos:7
MAINTAINER Tim Marcinowski <marshyski@gmail.com>

USER root

RUN yum -y update
RUN yum install -y epel-release gcc createrepo python-setuptools python-devel
RUN yum install -y python-pip supervisor
RUN pip install --upgrade pip
ADD ./requirements.txt /tmp/requirements.txt
RUN pip install -r /tmp/requirements.txt
RUN yum remove -y gcc && yum clean all
ADD . /opt/yum-nginx-api
ADD supervisor/yumapi.ini /etc/supervisord.d/yumapi.ini

EXPOSE 8888

CMD ["supervisord", "-n", "-c", "/etc/supervisord.conf"]
