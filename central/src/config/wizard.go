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

package config

import (
	"fmt"
	"github.com/SOF3/redox/central/src/util"
	"math"
	"strconv"
	"strings"
)

func runInstallWizard() (err error) {
	title("Redox Central setup wizard")
	info(`Redox Central is the heart of your server network. All servers will connect to it for communication.`)
	info(`Reminder: If your servers are hosted on separate machines, host Central on its own machine and keep its IP secret to protect it from attackers. If your servers are hosted on the same machine as Redox Central, remember not to open Central's port externally.`)

	Config.RedoxServer.Ip = "0.0.0.0"

	port, err := queryInt("Port of Central server", 0, 65535, 29456)
	if err != nil {
		return
	}
	Config.RedoxServer.Port = uint16(port)

	randomPassword, err := util.CryptoSecureRandomString(util.AlphaNum, 16)
	if err != nil {
		return
	}
	Config.RedoxServer.Password, err = queryString("Password of Central server", randomPassword)
	if err != nil {
		return
	}
	// querySSL(*Config.RedoxServer.SSL)
	Config.RedoxServer.SSL = nil

	title("Control Panel")
	info(`You can setup Redox Central in the future using a browser`)

	Config.ControlPanel.Ip = "0.0.0.0"

	port, err = queryInt("Port of Control Panel", 0, 65535, 29457)
	if err != nil {
		return
	}
	Config.ControlPanel.Port = uint16(port)

	randomPassword, err = util.CryptoSecureRandomString(util.AlphaNum, 16)
	if err != nil {
		return
	}
	Config.ControlPanel.Password, err = queryString("Password to login into Control Panel", randomPassword)
	// querySSL(*Config.ControlPanel.SSL)
	Config.ControlPanel.SSL = nil

	ip, err := util.GetMyIP()
	if err != nil {
		return
	}
	info("Setup complete! You can access the control panel through https://" + ip + ":29457")

	return SaveConfig()
}

func title(message string) { fmt.Println("=== " + message + " ===") }
func info(message string)  { fmt.Println("[*] " + message) }
func fatal(message string) { fmt.Println("!!! Fatal: " + message + " !!!") }
func warn(message string)  { fmt.Println("[!] " + message) }

// func querySSL(name string, field **TypeSSL) (err error) {
// 	enable, err := queryBoolean("Do you want to enable SSL for "+name+"?", true)
// 	if err != nil {
// 		return
// 	}
// 	if !enable {
// 		*field = nil
// 		return
// 	}
// 	ssl := new(TypeSSL)
// 	*field = ssl
//
// 	use, err := queryBoolean("Do you already have an SSL certificate? Type \"n\" to generate one", false)
// 	if err != nil {
// 		return
// 	}
// 	if use {
// 		for {
// 			ssl.CertFile, err = query("Path to certificate file", "", "", nil)
// 			if err != nil {
// 				return
// 			}
// 			if exists, err := util.FileExists(ssl.CertFile); err != nil || !exists {
// 				warn("Invalid file")
// 				use, err = queryBoolean("Do you really want to use an existing SSL certificate? Type \"n\" if you want to generate one.", false)
// 				if err != nil {
// 					return
// 				}
// 				if use {
// 					continue
// 				} else {
// 					goto gen
// 				}
// 			}
// 		}
// 		return
// 	}
//
// gen:
// 	return
// }

func queryBoolean(message string, defaultValue bool) (result bool, err error) {
	defaultString := util.QString(defaultValue, "y", "n")

	var char string
	char, err = query(message, "y/n, default: "+defaultString, defaultString, func(str string) bool {
		return strings.ToLower(str) == "y" || strings.ToLower(str) == "n"
	})
	if char == "y" {
		return true, nil
	}
	if char == "n" {
		return false, nil
	}
	panic("Program logic error")
}
func queryInt(message string, min int, max int, defaultValue int) (result int, err error) {
	var hint string
	var rangeChecker func(i int) bool
	if min == math.MinInt64 {
		// has no min range
		if max == math.MaxInt64 {
			// has no range
			hint = fmt.Sprintf("default %d", defaultValue)
			rangeChecker = nil
		} else {
			// has max range only
			hint = fmt.Sprintf("at most %d, default %d", max, defaultValue)
			rangeChecker = func(i int) bool { return i <= max }
		}
	} else {
		// has min range
		if max == math.MaxInt64 {
			// has min range only
			hint = fmt.Sprintf("at least %d, default %d", min, defaultValue)
			rangeChecker = func(i int) bool { return min <= i }
		} else {
			// has min and max range
			hint = fmt.Sprintf("%d-%d, default %d", min, max, defaultValue)
			rangeChecker = func(i int) bool { return min <= i && i <= max }
		}
	}

	var validator func(str string) bool
	if rangeChecker != nil {
		validator = func(str string) bool {
			i, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return false
			}
			return rangeChecker(int(i))
		}
	}

	s, err := query(message, hint, strconv.Itoa(defaultValue), validator)
	if err != nil {
		return
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i), nil
}
func queryString(message string, defaultValue string) (string, error) {
	return query(message, util.QString(defaultValue != "", "default: "+defaultValue, ""), defaultValue, nil)
}

func query(message string, hint string, defaultValue string, validator func(str string) bool) (result string, err error) {
	for {
		_, err = fmt.Print(util.QString(hint != "", "[?] "+message+" ("+hint+"): ", "[?] "+message+": "))
		if err != nil {
			return
		}
		buffer := ""
		isPrefix := true
		for isPrefix {
			var bytes []byte
			bytes, isPrefix, err = StdinReader.ReadLine()
			if err != nil {
				return
			}
			buffer += string(bytes)
		}
		str := string(buffer)

		if str == "" {
			result = defaultValue
		} else {
			result = str
		}

		if validator == nil || validator(result) {
			return
		}
		warn("Invalid input")
	}
}
