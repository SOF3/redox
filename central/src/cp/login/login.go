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

package login

import (
	"github.com/SOF3/redox/central/src/config"
	"github.com/SOF3/redox/central/src/cp/http_util"
	"github.com/SOF3/redox/central/src/util"
	"net/http"
	"strings"
)

type Middleware struct {
	sessions *util.ExpiringSyncMap
}

func New() *Middleware {
	return &Middleware{
		sessions: util.NewExpiringSyncMap(config.Config.ControlPanel.SessionDuration),
	}
}

func (m *Middleware) W(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		cookieObject, err := req.Cookie("RedoxCPSession")
		if err != nil && err != http.ErrNoCookie {
			http_util.HandleError(http_util.InternalError(err), res)
			return
		}

		var cookie string
		if err == http.ErrNoCookie {
			headers, exists := req.Header["Authorization"]

			var auth bool
			if exists && len(headers) > 0 {
				header := headers[0]
				if strings.HasPrefix(header, "Basic ") && header[len("Basic "):] == config.Config.ControlPanel.Password {
					auth = true
					cookie, err = m.sessions.FillRandom(util.AlphaNum, 32, nil)
					http.SetCookie(res, &http.Cookie{
						Name:     "RedoxCPSession",
						Path:     "/",
						Value:    cookie,
						HttpOnly: true,
						Secure:   config.Config.ControlPanel.SSL != nil,
						MaxAge:   int(config.Config.ControlPanel.SessionDuration.Seconds()),
					})
				}
			}
			if !auth {
				res.Header().Set("WWW-Authenticate", `Basic realm="Control panel password"`)
				res.WriteHeader(401)
				res.Write([]byte("Please login"))
				return
			}
		}

		cookie = cookieObject.Value
		_, exists := m.sessions.GetExists(cookie)

		if exists {
			handler(res, req)
		}
	}
}
