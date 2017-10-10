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
	"time"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	"github.com/h2non/filetype"
)

// crRoutine is a simple buffer to not overload the system
// by running too many createrepo system commands at the
// same time.
func crRoutine() {
	for {
		if crCtr > 0 {
			for i := 0; i < maxRetries; i++ {
				crExec := exec.Command("createrepo", "-v", "-p", "--update", "--workers", createRepo, uploadDir)
				_, err := crExec.Output()
				if err != nil {
					log.Println("Unable to execute createrepo ", err)
				} else {
					crCtr--
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
		time.Sleep(3 * time.Second)
	}
}

// uploadRoute function for handler /upload
func uploadRoute(c *routing.Context) error {
	c.Request.ParseMultipartForm(maxLength)
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.Response.WriteHeader(http.StatusInternalServerError)
		c.Write("Upload Failed")
		return err
	}
	defer file.Close()
	filePath := uploadDir + handler.Filename
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		c.Response.WriteHeader(http.StatusInternalServerError)
		c.Write("Upload Failed")
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	buf, _ := ioutil.ReadFile(filePath)
	if kind, err := filetype.Match(buf); err != nil || kind.MIME.Value != "application/x-rpm" {
		err := os.Remove(filePath)
		if err != nil {
			log.Println("Unable to delete " + filePath)
		}
		c.Response.WriteHeader(http.StatusUnsupportedMediaType)
		return c.Write(handler.Filename + " not RPM")
	}
	// If not in Development mode increment create-repo counter
	// for command to be ran by go routine crRoutine
	if !devMode {
		crCtr++
	}
	c.Response.WriteHeader(http.StatusAccepted)
	return c.Write("Uploaded")
}

// healthRoute function for handler /health
func healthRoute(c *routing.Context) error {
	c.Response.Header().Add("Version", commitHash)
	return c.Write("OK")
}

func main() {
	configValidate()
	go crRoutine()
	rtr := routing.New()

	api := rtr.Group("/api",
		fault.Recovery(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		content.TypeNegotiator(content.JSON))

	// Disable logging on health and metrics endpoints
	api.Get("/health", healthRoute)
	api.Use(access.Logger(log.Printf))
	api.Post("/upload", uploadRoute)

	http.Handle("/", rtr)
	log.Printf("yumapi built-on %s, version %s started on %s \n", builtOn, commitHash, port)
	http.ListenAndServe(":"+port, nil)
}
