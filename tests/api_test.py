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

import requests

baseurl = 'localhost'

''' Post RPM to API end point '''
posturl = 'http://%s/api/upload' % baseurl
files = {'file': open('yum-nginx-api-test-0.1-1.x86_64.rpm', 'rb')}
post = requests.post(posturl, files=files)
if post.status_code == 200:
  print "SAT: RPM post successful"
else:
  print "UNSAT: RPM post"

''' Verify 405 message '''
if requests.get(posturl).status_code == 405:
  print "SAT: 405 message recieved"
else:
  print "UNSAT: 405 message"

''' Verify 400 message '''
posturl = 'http://%s/api' % baseurl
if requests.get(posturl).status_code == 404:
  print "SAT: 400 message recieved"
else:
  print "UNSAT: 400 message"

''' Verify repo list '''
posturl = 'http://%s/api/repo' % baseurl
if requests.get(posturl).status_code == 200:
  print "SAT: repo list worked"
else:
  print "UNSAT: repo list didn't work"

''' Verify health check message '''
posturl = 'http://%s/api/health' % baseurl
if requests.get(posturl).status_code == 200:
  print "SAT: health check worked"
else:
  print "UNSAT: health check didn't work"
