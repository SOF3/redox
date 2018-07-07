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
	"bufio"
	"encoding/json"
	"github.com/SOF3/redox/central/src/util"
	"os"
)

const configFile = "data/config.json"

var Config Type

var StdinReader = bufio.NewReader(os.Stdin)

func Init() (err error) {
	util.SafeMkdir("data", false)
	exists, err := util.FileExists(configFile)
	if err != nil {
		return
	}

	if exists {
		var file *os.File
		file, err = os.Open(configFile)
		dec := json.NewDecoder(file)
		dec.DisallowUnknownFields()
		err = dec.Decode(&Config)
	} else {
		err = runInstallWizard()
	}
	return
}

func SaveConfig() (err error) {
	file, err := os.Create(configFile)
	if err != nil {
		return
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	return enc.Encode(Config)
}
