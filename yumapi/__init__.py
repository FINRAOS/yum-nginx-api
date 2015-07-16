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

from datetime import datetime, timedelta
from psutil import cpu_percent, process_iter
from subprocess import call
from flask import Flask, request, jsonify, abort
from flask_limiter import Limiter
from werkzeug.utils import secure_filename
from werkzeug.contrib.fixers import ProxyFix
import os
import magic
import repotojson
import yaml
import time

config_file = os.getcwd() + '/yumapi/config.yaml'
config_yaml = yaml.load(file(config_file, 'r'))
upload_dir = config_yaml['upload_dir']

"""Verify upload directory is set in config.yaml"""
if os.path.isdir(upload_dir) == False:
    print upload_dir, "doesn't exist, please create directory set directory in configurations.py"
    exit()

start_time = time.time()
hostname = os.uname()[1]
allowed_ext = set(['rpm'])
allowed_mime = set(['application/x-rpm'])

if config_yaml['request_limit']:
    request_limit = config_yaml['request_limit']
else:
    request_limit = '1 per second'

if config_yaml['createrepo_workers']:
    createrepo_workers = config_yaml['createrepo_workers']
else:
    createrepo_workers = '2'

if config_yaml['max_content_length']:
    max_content_length = config_yaml['max_content_length']
else:
    max_content_length = '90 * 1024 * 1024 * 1024'

app = Flask(__name__, static_folder='', static_url_path='')
app.config['upload_dir'] = upload_dir
"""900MB default limit set below"""
app.config['MAX_CONTENT_LENGTH'] = max_content_length

limiter = Limiter(app)

def uptime():
    """Return uptime about nginxify for health check"""
    seconds = timedelta(seconds=int(time.time() - start_time))
    d = datetime(1,1,1) + seconds
    return("%dD:%dH:%dM:%dS" % (d.day-1, d.hour, d.minute, d.second))

def allowed_file(filename):
   """Verify filename is in allowed extension"""
   return '.' in filename and \
      filename.rsplit('.', 1)[1] in allowed_ext

@app.route('/api/upload', methods=['POST'])
@limiter.limit(request_limit)
def upload_file():
   """Post RPMs to upload URI, this does all the work"""
   file = request.files['file']
   filename = secure_filename(file.filename)
   if file and allowed_file(file.filename):
      file.save(os.path.join(app.config['upload_dir'], filename))
      file2 = upload_dir + '/' + filename
      file_mime = magic.from_file(file2, mime=True)
      filesize = os.path.getsize(file2) >> 20
      if file_mime in allowed_mime and os.path.getsize(file2) >= 1870:
         call(["createrepo", "-v", "-p", "--update", "--workers", createrepo_workers, upload_dir])
	 return jsonify(name=filename, size_mb=filesize, mime=file_mime, status=202)
      else:
	 os.remove(file2)
	 abort(415)
   else:
      abort(415)


@app.route('/api/repo')
@limiter.limit(request_limit)
def list_repo():
   """Return json response from sqlite database yum has"""
   repotojson.main()
   return app.send_static_file('repo.json')

@app.route('/api/health')
@limiter.limit(request_limit)
def health():
    """Health URI to verify yum-nginx-api is operational"""
    return jsonify(hostname=hostname, uptime=uptime(), cpu_percent=int(cpu_percent(interval=None, percpu=False)), status=200)

"""Error handlers for main error responses"""
@app.errorhandler(404)
def page_not_found(error):
   return jsonify(message='Invaild URI, did you mean /api/upload, /api/repo or /api/health', status=404), 404

@app.errorhandler(405)
def method_not_allowed(error):
   return jsonify(message='Method not allowed, post binary RPM to URI', status=405), 405

@app.errorhandler(415)
def unsupported_media(error):
   return jsonify(message='File not RPM', status=415), 415

app.wsgi_app = ProxyFix(app.wsgi_app)

if __name__ == '__main__':
     app.run('127.0.0.1', debug = False)
