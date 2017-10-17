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

package main

import (
	"errors"
	"os"
	"path"

	"github.com/FINRAOS/yum-nginx-api/repojson"
	"github.com/spf13/viper"
)

const (
	crError = "config: createrepo_workers is less than 1"
	mlError = "config: max_content_length is less than 1MB"
	upError = "config: upload_directory does not exist"
	ptError = "config: port is not above port 80"
	mxError = "config: max_retries is less than 1"
)

// Some vars are used with different types in handlers vs validation
var (
	builtOn    string
	commitHash string
	createRepo string
	maxLength  int64
	uploadDir  string
	port       string
	devMode    bool
	maxRetries int
	crCtr      int64
	crPaths    = [2]string{"/bin/createrepo", "/usr/bin/createrepo"}
	rJSON      []repojson.Repo
	crBin      string
)

// Validate configurations and if createrepo binary is present in path
func configValidate() error {
	viper.SetConfigName("yumapi")
	viper.AddConfigPath("/opt/yum-nginx-api/yumapi/")
	viper.AddConfigPath("/etc/yumapi/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.New("Fatal error config file: " + err.Error())
	}

	viper.SetDefault("createrepo_workers", 1)
	viper.SetDefault("max_content_length", 10000000)
	viper.SetDefault("upload_dir", "./")
	viper.SetDefault("port", 8080)
	viper.SetDefault("dev_mode", false)
	viper.SetDefault("max_retries", 3)

	createRepo = viper.GetString("createrepo_workers")
	maxLength = viper.GetInt64("max_content_length")
	uploadDir = path.Clean(viper.GetString("upload_dir")) + "/"
	port = viper.GetString("port")
	devMode = viper.GetBool("dev_mode")
	maxRetries = viper.GetInt("max_retries")

	if viper.GetInt64("createrepo_workers") < 1 {
		return errors.New(crError)
	}
	if maxLength < 1000000 {
		return errors.New(mlError)
	}
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		return errors.New(upError)
	}
	if viper.GetInt64("port") < 80 {
		return errors.New(ptError)
	}
	if maxRetries < 1 {
		return errors.New(mxError)
	}
	if !devMode {
		for _, cr := range crPaths {
			if _, err := os.Stat(cr); !os.IsNotExist(err) {
				crBin = cr
				break
			}
		}
		if crBin == "" {
			return errors.New(crError)
		}
	}
	return nil
}
