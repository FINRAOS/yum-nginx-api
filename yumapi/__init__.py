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
from flask import Flask, request, jsonify
from werkzeug.utils import secure_filename
from werkzeug.contrib.fixers import ProxyFix
from subprocess import call
import configuration as config
import repotojson

upload_dir    = 'config.upload_dir'
allowed_ext   = set(['rpm'])
allowed_mime  = set(['application/x-rpm'])

app = Flask(__name__, static_folder='', static_url_path='')
app.config['upload_dir'] = upload_dir
#900MB limit set below
app.config['MAX_CONTENT_LENGTH'] = 90 * 1024 * 1024 * 1024

def allowed_file(filename):
    return '.' in filename and \
           filename.rsplit('.', 1)[1] in allowed_ext

@app.errorhandler(404)
def page_not_found(e):
    return jsonify(message='Invaild URL, did you mean /api/upload or /api/repo', status=404), 404

@app.errorhandler(405)
def method_not_allowed(e):
    return jsonify(message='Method not allowed, post binary RPM to URL', status=405), 405

@app.route('/api/upload', methods=['POST'])
def upload_file():
    file = request.files['file']
    filename = secure_filename(file.filename)
    if file and allowed_file(file.filename):
       file.save(os.path.join(app.config['upload_dir'], filename))
       file2 = upload_dir + '/' + filename
       file_mime = magic.from_file(file2, mime=True)
       filesize = os.path.getsize(file2) >> 20
       if file_mime in allowed_mime:
          call(["createrepo", "-v", "-p", "--update", "--workers", "2", upload_dir])
	  return jsonify(name=filename, size_mb=filesize, mime=file_mime, status=202)
       else:
	  os.remove(file2)
	  return jsonify(name=filename, size_mb=filesize, mime=file_mime, status=415)
    else:
       return jsonify(name=filename, status=415)

@app.route('/api/repo')
def list_repo():
    repotojson.main()
    return app.send_static_file('repo.json')

app.wsgi_app = ProxyFix(app.wsgi_app)

if __name__ == '__main__':
     app.run('0.0.0.0', debug = True)
