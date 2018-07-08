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

package util

import (
	"encoding/json"
	"net/http"
)

func GetMyIP() (ip string, err error) {
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		return
	}

	dec := json.NewDecoder(resp.Body)
	var m = map[string]interface{}{}
	err = dec.Decode(&m)
	if err != nil {
		return
	}
	return m["ip"].(string), nil
}
