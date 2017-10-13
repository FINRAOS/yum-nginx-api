#!/bin/bash
#  (C) Copyright 2014 yum-nginx-api Contributors.
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at 
#  http://www.apache.org/licenses/LICENSE-2.0
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

cat << 'LIMITS' > /etc/security/limits.d/nginx-limits.conf
nginx       soft    nofile  30000
nginx       hard    nofile  65536
nginx       soft    nproc   16384
nginx       hard    nproc   16384
nginx       hard    stack   1024
LIMITS

cat << 'SYSCTL' >> /etc/sysctl.d/nginx.conf
fs.file-max=70000
fs.nr_open=70000
net.core.netdev_max_backlog=4096
net.core.rmem_max=16777216
net.core.somaxconn=65535
net.core.wmem_max=16777216
net.ipv4.tcp_fin_timeout=30
net.ipv4.tcp_keepalive_time=30
net.ipv4.tcp_max_syn_backlog=20480
net.ipv4.tcp_max_tw_buckets=400000
net.ipv4.tcp_no_metrics_save=1
net.ipv4.tcp_syn_retries=2
net.ipv4.tcp_synack_retries=2
net.ipv4.tcp_tw_recycle=1
net.ipv4.tcp_tw_reuse=1
vm.min_free_kbytes=65536
vm.overcommit_memory=1
SYSCTL
