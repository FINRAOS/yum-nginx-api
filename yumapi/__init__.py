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

import os
import magic
from configuration import upload_dir
import repotojson
from healthcheck import HealthCheck
from subprocess import call
from flask import Flask, request, jsonify
from werkzeug.utils import secure_filename
from werkzeug.contrib.fixers import ProxyFix

allowed_ext   = set(['rpm'])
allowed_mime  = set(['application/x-rpm'])

''' Verify upload directory is set in configuration.py '''
if os.path.isdir(upload_dir) == False:
    print upload_dir, "doesn't exist, please create directory set directory in configurations.py"
    exit()

app = Flask(__name__, static_folder='', static_url_path='')
app.config['upload_dir'] = upload_dir
#900MB limit set below
app.config['MAX_CONTENT_LENGTH'] = 90 * 1024 * 1024 * 1024

def allowed_file(filename):
   ''' Verify filename is in allowed extension '''
   return '.' in filename and \
      filename.rsplit('.', 1)[1] in allowed_ext

@app.route('/api/upload', methods=['POST'])
def upload_file():
   ''' Post RPMs to upload URI, this does all the work '''
   file = request.files['file']
   filename = secure_filename(file.filename)
   if file and allowed_file(file.filename):
      file.save(os.path.join(app.config['upload_dir'], filename))
      file2 = upload_dir + '/' + filename
      file_mime = magic.from_file(file2, mime=True)
      filesize = os.path.getsize(file2) >> 20
      if file_mime in allowed_mime and os.path.getsize(file2) >= 1870:
         call(["createrepo", "-v", "-p", "--update", "--workers", "2", upload_dir])
	 return jsonify(name=filename, size_mb=filesize, mime=file_mime, status=202)
      else:
	 os.remove(file2)
	 return jsonify(name=filename, size_mb=filesize, mime=file_mime, status=415)
   else:
      return jsonify(name=filename, status=415)

@app.route('/api/repo')
def list_repo():
   ''' Return json response from sqlite database yum has '''
   os.remove('repo.json')
   repotojson.main()
   return app.send_static_file('repo.json')

''' Health URI to verify yum-nginx-api is operational '''
health = HealthCheck(app, "/api/health")

''' Error handlers for main error reponses '''
@app.errorhandler(404)
def page_not_found(e):
   return jsonify(message='Invaild URI, did you mean /api/upload, /api/repo or /api/health', status=404), 404

@app.errorhandler(405)
def method_not_allowed(e):
   return jsonify(message='Method not allowed, post binary RPM to URI', status=405), 405

app.wsgi_app = ProxyFix(app.wsgi_app)

if __name__ == '__main__':
     app.run('0.0.0.0', debug = False)
