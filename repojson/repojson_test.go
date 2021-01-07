package repojson

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const jstring = `[{"name":"yum-nginx-api-test","arch":"x86_64","version":"0.1","summary":"Yum NGINX API Test RPM"}]`

// TestSetPsqlite finds primary.sqlite archive in current directory
func TestSetPsqlite(t *testing.T) {
	var pSqlite string
	err := filepath.Walk("./", setPsqlite(&pSqlite))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if pSqlite == "" {
		t.Fatalf("Expected to match file")
	}
}

// TestRegExSetPsqlite tests regex used to find files in given path
func TestRegExSetPsqlite(t *testing.T) {
	var pSqlite string
	b := []byte("")
	if err := ioutil.WriteFile("/tmp/primary.sqlitexz", b, 0644); err != nil {
		os.Remove("/tmp/primary.sqlitexz")
		t.Fatalf("Error writing file")
	}
	err := filepath.Walk("/tmp", setPsqlite(&pSqlite))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if pSqlite != "" {
		os.Remove("/tmp/primary.sqlitexz")
		t.Fatalf("Expected to not match file got " + pSqlite)
	}
	os.Remove("/tmp/primary.sqlitexz")

	if err := ioutil.WriteFile("/tmp/primary.sqlite.xzz", b, 0644); err != nil {
		os.Remove("/tmp/primary.sqlite.xzz")
		t.Fatalf("Error writing file")
	}
	err = filepath.Walk("/tmp", setPsqlite(&pSqlite))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if pSqlite != "" {
		os.Remove("/tmp/primary.sqlite.xzz")
		t.Fatalf("Expected to not match file got " + pSqlite)
	}
	os.Remove("/tmp/primary.sqlite.xzz")

	if err := ioutil.WriteFile("/tmp/primary.sqlite.bz", b, 0644); err != nil {
		os.Remove("/tmp/primary.sqlite.bz")
		t.Fatalf("Error writing file")
	}
	err = filepath.Walk("/tmp", setPsqlite(&pSqlite))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if pSqlite != "" {
		os.Remove("/tmp/primary.sqlite.bz")
		t.Fatalf("Expected to not match file got " + pSqlite)
	}
	os.Remove("/tmp/primary.sqlite.bz")
}

// TestExtractBZ extracts bzip2 archive and verifies uncompressed
func TestExtractBZ(t *testing.T) {
	err := extractBZ2("./", "./primary.sqlite.bz2")
	if err != nil {
		t.Fatalf("Expected to extract test file")
	}
	if _, err := os.Stat("./primary.sqlite"); os.IsNotExist(err) {
		t.Fatalf("Expected extracted file to exist")
	}
	os.Remove("./primary.sqlite")
}

// TestExtractXZ extracts xz archive and verifies uncompressed

func TestExtractXZ(t *testing.T) {
	err := extractXZ("./", "./primary.sqlite.xz")
	if err != nil {
		t.Fatalf("Expected to extract test file")
	}
	if _, err := os.Stat("./primary.sqlite"); os.IsNotExist(err) {
		t.Fatalf("Expected extracted file to exist")
	}
	os.Remove("./primary.sqlite")
}

// TestRepoJSON matches JSON constant to archive
// after it is read from sqlite and marshalled
func TestRepoJSON(t *testing.T) {
	var ba []Repo
	ba, err := RepoJSON("./")
	if err != nil {
		os.Remove("primary.sqlite")
		t.Fatalf("Expected to find primary.sqlite")
	}
	js, _ := json.Marshal(ba)
	if string(js) != jstring {
		t.Fatalf("Expect %s got %s", jstring, string(js))
	}
}

// TestFailRepoJSON verifies RepoJSON return errors if no match is made
// in directory path
func TestFailRepoJSON(t *testing.T) {
	dir, err := ioutil.TempDir("", "yumapi")
	if err != nil {
		t.Fatalf("Failed to make temp dir")
	}
	if _, err := RepoJSON(dir); err == nil {
		os.Remove("primary.sqlite")
		os.RemoveAll(dir)
		t.Fatalf("Expected to not find primary.sqlite")
	}
	os.Remove("primary.sqlite")
	os.RemoveAll(dir)
}
