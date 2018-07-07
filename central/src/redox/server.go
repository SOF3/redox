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

package redox

import (
	"github.com/SOF3/redox/central/src/config"
	"github.com/SOF3/redox/central/src/redox/node"
	"github.com/SOF3/redox/central/src/util/ch"
	"log"
	"net"
	"strconv"
)

type Server struct {
	socket      net.Listener
	clientsChan chan net.Conn
	errorsChan  chan error
}

func (server *Server) Run(redox2main chan<- error, main2redox <-chan int) {
	var err error
	cfg := config.Config.RedoxServer

	addr := cfg.Ip + ":" + strconv.Itoa(int(cfg.Port))
	server.socket, err = net.Listen("tcp", addr)
	if err != nil {
		redox2main <- err
		return
	}
	log.Println("[Redox Server] Listening on " + addr)

	go server.acceptClients()

	for {
		select {
		case client := <-server.clientsChan:
			log.Println("Connection opened from " + client.RemoteAddr().String())
			node.New(client)
		case signal := <-main2redox:
			log.Println("Caught signal " + strconv.Itoa(signal))
			switch signal {
			case ch.Quit:
				log.Println("[Redox Server] Stopping...")
				err = server.socket.Close()
				if err != nil {
					panic(err)
				}
				close(redox2main)
				return
			default:
				log.Println("Unknown main goroutine signal " + strconv.Itoa(signal))
			}
		}
	}
}

func (server *Server) acceptClients() {
	for {
		client, err := server.socket.Accept()
		if err != nil {
			server.errorsChan <- err
			close(server.clientsChan)
			return
		}
		server.clientsChan <- client
	}
}
