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

package main

import (
	"github.com/SOF3/redox/central/src/config"
	"github.com/SOF3/redox/central/src/cp"
	"github.com/SOF3/redox/central/src/redox"
	"github.com/SOF3/redox/central/src/util/ch"
	"log"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	log.Println("Starting redox server")
	redox2main := make(chan error, 1)
	main2redox := make(chan int, 1)
	go new(redox.Server).Run(redox2main, main2redox)

	log.Println("Starting control panel")
	cp2main := make(chan error, 1)
	main2cp := make(chan int, 1)
	go new(cp.Server).Run(cp2main, main2cp)

	os2main := make(chan os.Signal, 1)
	signal.Notify(os2main, os.Interrupt)

	log.Println("Process ID: " + strconv.Itoa(os.Getpid()))
	log.Println("Press Ctrl-C to stop Redox Central")
	select {
	case err = <-redox2main:
		log.Fatalln(err.Error())
		main2cp <- ch.Quit
		panic(err)
	case err = <-cp2main:
		log.Fatalln(err.Error())
		main2redox <- ch.Quit
		panic(err)
	case <-os2main:
		log.Println("Received sigint")
		main2redox <- ch.Quit
		main2cp <- ch.Quit
		<-redox2main
		<-cp2main
		log.Println("Redox Central has been stopped.")
	}

}
