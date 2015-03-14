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
#  Main contributor is Pierre-Yves Chibon https://github.com/pypingou

import pkg_resources
import contextlib
import json
import os
import shutil
import sys
import tempfile
from sqlalchemy import Column, ForeignKey, Integer, Text, create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from configuration import upload_dir

BASE = declarative_base()

class Package(BASE):
    ''' Maps the packages table in the primary.sqlite database from
    repodata to a python object.
    '''
    __tablename__ = 'packages'
    pkgKey = Column(Integer, primary_key=True)
    name = Column(Text)
    rpm_sourcerpm = Column(Text)
    version = Column(Text)
    arch = Column(Text)
    summary = Column(Text)

    @property
    def basename(self):
        ''' Return the base package name using the rpm_sourcerpms info. '''
        return self.rpm_sourcerpm.rsplit('-', 2)[0]


def find_primary_sqlite(paths):
    ''' Find all the primary.sqlite files located at or under the given path.'''
    if not isinstance(paths, list):
        paths = [paths]
    files = []
    for path in paths:
        if not os.path.isdir(path):
            continue
        for (dirpath, dirnames, filenames) in os.walk(path):
            for filename in filenames:
                if 'primary.sqlite' in filename:
                    files.append(os.path.join(dirpath, filename))
    return files


def decompress_primary_db(archive, location):
    ''' Decompress the given XZ archive at the specified location. '''
    if archive.endswith('.xz'):
        import lzma
        with contextlib.closing(lzma.LZMAFile(archive)) as stream_xz:
            data = stream_xz.read()
        with open(location, 'wb') as stream:
            stream.write(data)
    elif archive.endswith('.gz'):
        import tarfile
        with tarfile.open(archive) as tar:
            tar.extractall(path=location)
    elif archive.endswith('.bz2'):
        import bz2
        with open(location, 'w') as out:
            bzar = bz2.BZ2File(archive)
            out.write(bzar.read())
            bzar.close()
    elif archive.endswith('.sqlite'):
        with open(location, 'w') as out:
            with open(archive) as inp:
                out.write(inp.read())


def get_pkg_info(session, pkg_name):
    ''' Query the sqlite database for the package specified. '''
    pkg = session.query(Package).filter(Package.name == pkg_name).one()
    return pkg


def main():
    working_dir = tempfile.mkdtemp(prefix='repotojson-')
    output = {}
    dbfiles = find_primary_sqlite(upload_dir)
    for dbfile_xz in dbfiles:
        cur_fold = os.path.join(*dbfile_xz.rsplit(os.sep, 2)[:-2])
        dbfile = os.path.join(working_dir, 'primary_db_.sqlite')
        decompress_primary_db(dbfile_xz, dbfile)
        if not os.path.isfile(dbfile):
            print '%s was incorrectly decompressed -- ignoring' % dbfile
            continue
        db_url = 'sqlite:///%s' % dbfile
        db_session = sessionmaker(bind=create_engine(db_url))
        session = db_session()
        cnt = 0
        new = 0
        for pkg in session.query(Package).all():
            if pkg.basename in output:
                if pkg.arch not in output[pkg.basename]['arch']:
                    output[pkg.basename]['arch'].append(pkg.arch)
            else:
                new += 1
                output[pkg.basename] = {
                  'arch': [pkg.arch],
                  'version': pkg.version,
                  'summary': pkg.summary,
                }
            cnt += 1
        outputfile = 'repo.json'
        with open(outputfile, 'w') as stream:
            stream.write(json.dumps(output, sort_keys=True, indent=2, separators=(',', ': ')))
    shutil.rmtree(working_dir)


if __name__ == '__main__':
    main()
