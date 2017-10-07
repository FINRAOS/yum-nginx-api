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

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	"github.com/h2non/filetype"
)

func main() {
	rtr := routing.New()
	api := rtr.Group("/api",
		fault.Recovery(log.Printf),
		slash.Remover(http.StatusMovedPermanently))

	// Do not enable logging on health and metrics endpoints
	api.Get("/health", func(c *routing.Context) error {
		return c.Write("OK")
	})

	rtr.Use(
		access.Logger(log.Printf),
	)

	api.Use(
		content.TypeNegotiator(content.JSON),
	)

	api.Post("/upload", func(c *routing.Context) error {
		c.Request.ParseMultipartForm(32 << 20)
		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			return err
		}
		defer file.Close()
		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0664)
		if err != nil {
			return err
		}
		defer f.Close()
		io.Copy(f, file)
		buf, _ := ioutil.ReadFile(handler.Filename)
		isRPM, err := filetype.Match(buf)
		if err != nil {
			os.Remove(handler.Filename)
			return c.Write(isRPM.MIME)
		}
		return c.Write(handler.Filename)
	})

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}
