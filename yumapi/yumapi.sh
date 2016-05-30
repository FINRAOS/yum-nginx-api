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

NAME=yumapi
USER=root
GROUP=root
WORKERS=`lscpu | grep ^'CPU(s)' | awk '{ print $2 }'`
DEPLOY_DIR=/opt/yum-nginx-api

if [[ $1 = "" ]]; then
   LISTEN_ADDR=127.0.0.1
else
   LISTEN_ADDR=0.0.0.0
fi

cd $DEPLOY_DIR

# Start your unicorn
gunicorn $NAME:app -b $LISTEN_ADDR:8888 \
  --name $NAME \
  --workers $WORKERS \
  --user=$USER --group=$GROUP
