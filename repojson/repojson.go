//  (C) Copyright 2014 yum-nginx-api Contributors.
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//  http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package repojson

import (
	"compress/bzip2"
	"database/sql"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"

	// Bug in mattn/go-sqlite3 and also lighter
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/ulikunitz/xz"
)

// Repo fields is column names from packages table
// in primary.sqlite file. Other columns not exposed yet
// pkgKey,pkgId,epoch,release,description,url,
// time_file,time_build,rpm_license,rpm_vendor,
// rpm_group,rpm_buildhost,rpm_sourcerpm,rpm_header_start,
// rpm_header_end,rpm_packager,size_package,size_installed,
// size_archive,location_href,location_base,checksum_type
type Repo struct {
	Name    string `json:"name"`
	Arch    string `json:"arch"`
	Version string `json:"version"`
	Summary string `json:"summary"`
}

// setPsqlite sets string pSqlite by walking
// path given to WalkFunc to find if primary.sqlite
// is present in directory with bzip2 or xz extenstion
func setPsqlite(pSqlite *string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		match, _ := regexp.MatchString("primary.sqlite.(xz|bz2)$", path)
		if match {
			*pSqlite = path
			return nil
		}
		return nil
	}
}

// extractXZ extracts primary.sqlite xz if exists
func extractXZ(path, pSqlite string) error {
	f, err := os.Open(pSqlite)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := xz.NewReader(f)
	if err != nil {
		return err
	}
	w, err := os.Create(path + "primary.sqlite")
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err = io.Copy(w, r); err != nil {
		return err
	}

	return nil
}

// extractBZ2 extracts primary.sqlite bzip2 if exists
func extractBZ2(path, pSqlite string) error {
	f, err := os.Open(pSqlite)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bzip2.NewReader(f)
	if err != nil {
		return err
	}
	w, err := os.Create(path + "primary.sqlite")
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err = io.Copy(w, r); err != nil {
		return err
	}

	return nil
}

// repoSqlite connects to uncompressed sqlite file
// quries database and returns array of Repo objects
func repoSqlite(path string) ([]Repo, error) {
	var repo []Repo
	db, err := sql.Open("sqlite3", path+"primary.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlStmt := `select name, arch, version, summary from packages;`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r := Repo{}
		err = rows.Scan(&r.Name, &r.Arch, &r.Version, &r.Summary)
		if err != nil {
			return nil, err
		}
		repo = append(repo, r)
	}
	return repo, nil
}

// RepoJSON is main function in repojson package to return
// array of Repo objects
func RepoJSON(path string) ([]Repo, error) {
	var (
		extract error
		pSqlite string
	)
	filepath.Walk(path, setPsqlite(&pSqlite))
	if pSqlite != "" {
		switch filepath.Ext(pSqlite) {
		case ".xz":
			extract = extractXZ(path, pSqlite)
		case ".bz2":
			extract = extractBZ2(path, pSqlite)
		}
		if extract == nil {
			j, err := repoSqlite(path)
			if err != nil {
				return nil, extract
			}
			return j, nil
		}
		return nil, extract
	}
	return nil, errors.New("RepoJSON: couldn't find a supported primary.sqlite file in " + path)
}
