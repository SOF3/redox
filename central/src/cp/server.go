/*
 * redox/central
 *
 * Copyright (C) 2018 SOFe
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cp

import (
	"github.com/SOF3/redox/central/src/config"
	"github.com/SOF3/redox/central/src/cp/login"
	"github.com/SOF3/redox/central/src/util/ch"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	server *http.Server

	login *login.Middleware
}

func (server *Server) Run(cp2main chan<- error, main2cp <-chan int) {
	cfg := config.Config.ControlPanel

	server.server = &http.Server{
		Addr:    cfg.Ip + ":" + strconv.Itoa(int(cfg.Port)),
		Handler: http.NewServeMux(),
	}
	server.login = login.New()
	server.setupRoutes()

	errorsChan := make(chan error, 1)
	go func() {
		var err error
		if cfg.SSL != nil {
			err = server.server.ListenAndServeTLS(cfg.SSL.CertFile, cfg.SSL.KeyFile)
		} else {
			err = server.server.ListenAndServe()
		}
		if err != nil {
			errorsChan <- err
		}
	}()

	log.Println("[Control panel] Listening on " + server.server.Addr)
	for {
		select {
		case err := <-errorsChan:
			cp2main <- err
			return
		case signal := <-main2cp:
			switch signal {
			case ch.Quit:
				err := server.server.Close()
				if err != nil {
					panic(err)
				}
				log.Println("Successfully shut down cp server")
				close(cp2main)
				return
			}
		}
	}
}

func (server *Server) handle(pattern string, handlerFunc http.HandlerFunc) {
	server.server.Handler.(*http.ServeMux).HandleFunc(pattern,
		server.login.W(handlerFunc))
}

func (server *Server) setupRoutes() {
	server.handle("/", func(writer http.ResponseWriter, request *http.Request) {
		http.NotFound(writer, request)
	})
}
