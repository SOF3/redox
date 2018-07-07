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

package http_util

import (
	"log"
	"net/http"
)

type Error struct {
	Internal bool
	Code     int
	Error    string
}

func InternalError(err error) Error {
	return Error{
		Internal: true, Error: err.Error(),
	}
}
func UserError(code int, err string) Error {
	return Error{
		Internal: false, Code: code, Error: err,
	}
}

func HandleError(err Error, res http.ResponseWriter) {
	if err.Internal {
		http.Error(res, "internal error", err.Code)
		log.Println("500 Internal Server Error: " + err.Error)
	} else {
		http.Error(res, err.Error, err.Code)
	}
}
