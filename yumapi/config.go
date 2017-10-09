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
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/viper"
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
)

func init() {
	viper.SetConfigName("yumapi")
	viper.AddConfigPath("/opt/yum-nginx-api/yumapi/")
	viper.AddConfigPath("/etc/yumapi/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.SetDefault("createrepo_workers", 2)
	viper.SetDefault("max_content_length", 10000000)
	viper.SetDefault("upload_dir", "./upload/")
	viper.SetDefault("port", 8080)
	viper.SetDefault("dev_mode", false)
}

func configValidate() {
	createRepo = viper.GetString("createrepo_workers")
	maxLength = viper.GetInt64("max_content_length")
	uploadDir = viper.GetString("upload_dir")
	port = viper.GetString("port")
	devMode = viper.GetBool("dev_mode")

	if viper.GetInt64("createrepo_workers") < 1 {
		panic(fmt.Errorf("createrepo_workers is less than 1"))
	}
	if maxLength < 1000000 {
		panic(fmt.Errorf("max_content_length is less than 1MB"))
	}
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		panic(fmt.Errorf("upload_directory %s does not exist", uploadDir))
	}
	if viper.GetInt64("port") < 80 {
		panic(fmt.Errorf("port is not above port 80"))
	}
	if !devMode {
		crBin := exec.Command("which", "createrepo")
		_, err := crBin.Output()
		if err != nil {
			panic(fmt.Errorf("createrepo binary not found in path"))
		}
	}
}
