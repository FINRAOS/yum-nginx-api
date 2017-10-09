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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	"github.com/h2non/filetype"
)

// Upload function for handler /upload
func upload(c *routing.Context) error {
	c.Request.ParseMultipartForm(maxLength)
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.Write("Upload failed")
		return err
	}
	defer file.Close()
	f, err := os.OpenFile(uploadDir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		c.Write("Upload failed")
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	buf, _ := ioutil.ReadFile(uploadDir + handler.Filename)
	_, err = filetype.Match(buf)
	if err != nil {
		err := os.Remove(uploadDir + handler.Filename)
		if err != nil {
			log.Println("Unable to delete " + uploadDir + handler.Filename)
		}
		return c.Write(handler.Filename + " not RPM")
	}
	if !devMode {
		crExec := exec.Command("createrepo", "-v", "-p", "--update", "--workers", createRepo, uploadDir)
		_, err := crExec.Output()
		if err != nil {
			log.Println("Unable to execute createrepo ", err)
			return c.Write(handler.Filename + " Unable to add to repo")
		}
	}
	return c.Write(handler.Filename)
}

// Health function for handler /health
func health(c *routing.Context) error {
	c.Response.Header().Add("Version", commitHash)
	return c.Write("OK")
}

func main() {
	configValidate()
	rtr := routing.New()
	api := rtr.Group("/api", fault.Recovery(log.Printf), slash.Remover(http.StatusMovedPermanently))

	// Do not enable logging on health and metrics endpoints
	api.Get("/health", health)

	rtr.Use(access.Logger(log.Printf))

	api.Use(content.TypeNegotiator(content.JSON))

	api.Post("/upload", upload)

	http.Handle("/", rtr)
	log.Printf("yumapi built-on %s, version %s started on %s \n", builtOn, commitHash, port)
	http.ListenAndServe(":"+port, nil)
}
