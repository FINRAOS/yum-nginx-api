#!/usr/bin/python
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

from locust import HttpLocust, TaskSet, task

class YumSvr(TaskSet):

  @task
  def index(self):
    self.client.get("/")

class MyLocust(HttpLocust):
  host = "http://localhost/"
  task_set = YumSvr
  min_wait = 5000
  max_wait = 15000
