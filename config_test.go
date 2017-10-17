package main

import (
	"io/ioutil"
	"os"
	"testing"
)

const testFile = "./yumapi.yaml"

// TestConfigValidate validates all configuration
// options in yumapi config file
func TestConfigValidate(t *testing.T) {
	b := []byte("upload_dir: .\ncreaterepo_workers: 2\ndev_mode: true\nmax_content_length: 220000000\nport: 80\nmax_retries: 1")
	if err := ioutil.WriteFile(testFile, b, 0644); err != nil {
		t.Fatalf("Error writing file")
	}
	if err := configValidate(); err != nil {
		os.Remove(testFile)
		t.Fatal(err)
	}
	os.Remove(testFile)
}

// TestConfigCR verifies error returned for zero workers
func TestConfigCR(t *testing.T) {
	b := []byte("createrepo_workers: 0")
	if err := ioutil.WriteFile(testFile, b, 0644); err != nil {
		t.Fatalf("Error writing file")
	}
	if err := configValidate(); err.Error() != crError {
		os.Remove(testFile)
		t.Fatal(err)
	}
	os.Remove(testFile)
}

// TestConfigML verifies max_content_length is not too low
func TestConfigML(t *testing.T) {
	b := []byte("max_content_length: 900000")
	if err := ioutil.WriteFile(testFile, b, 0644); err != nil {
		t.Fatalf("Error writing file")
	}
	if err := configValidate(); err.Error() != mlError {
		os.Remove(testFile)
		t.Fatal(err)
	}
	os.Remove(testFile)
}

// TestConfigUp verifies upload_dir does not exist
func TestConfigUP(t *testing.T) {
	b := []byte("upload_dir: /supp")
	if err := ioutil.WriteFile(testFile, b, 0644); err != nil {
		t.Fatalf("Error writing file")
	}
	if err := configValidate(); err.Error() != upError {
		os.Remove(testFile)
		t.Fatal(err)
	}
	os.Remove(testFile)
}

// TestConfigPT verifies ports under 80 are not allowed
func TestConfigPT(t *testing.T) {
	b := []byte("port: 22")
	if err := ioutil.WriteFile(testFile, b, 0644); err != nil {
		t.Fatalf("Error writing file")
	}
	if err := configValidate(); err.Error() != ptError {
		os.Remove(testFile)
		t.Fatal(err)
	}
	os.Remove(testFile)
}

// TestConfigMX verifies max retries are not less than 1
func TestConfigMX(t *testing.T) {
	b := []byte("max_retries: 0")
	if err := ioutil.WriteFile(testFile, b, 0644); err != nil {
		t.Fatalf("Error writing file")
	}
	if err := configValidate(); err.Error() != mxError {
		os.Remove(testFile)
		t.Fatal(err)
	}
	os.Remove(testFile)
}
