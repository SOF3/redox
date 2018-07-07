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
	"encoding/base64"
	"github.com/SOF3/redox/central/src/config"
	"github.com/SOF3/redox/central/src/cp/http_util"
	"github.com/SOF3/redox/central/src/util"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Middleware struct {
	sessions *util.ExpiringSyncMap

	bans *util.ExpiringSyncMap

	failures *util.ExpiringSyncMap
}

func New() *Middleware {
	return &Middleware{
		sessions: util.NewExpiringSyncMap(config.Config.ControlPanel.SessionDuration),
		bans:     util.NewExpiringSyncMap(util.Duration{Duration: time.Hour}),
		failures: util.NewExpiringSyncMap(util.Duration{Duration: time.Hour}),
	}
}

type failures struct {
	mutex sync.Mutex
	cnt   int
}

func (m *Middleware) W(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		host, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			http_util.HandleError(http_util.InternalError(err), res)
		}

		if _, exists := m.bans.GetExists(host); exists {
			http_util.HandleError(http_util.UserError(403, "Too many failed logins"), res)
			return
		}

		cookieObject, err := req.Cookie("RedoxCPSession")
		if err != nil && err != http.ErrNoCookie {
			http_util.HandleError(http_util.InternalError(err), res)
			return
		}

		var cookie string
		if err == nil {
			cookie = cookieObject.Value
			_, exists := m.sessions.GetExists(cookie)
			if exists {
				handler(res, req)
				return
			}
		}
		header := req.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Basic ") {
			goto invalidRequest
		}

		{
			passwordBytes, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(header, "Basic "))
			password := string(passwordBytes)
			index := strings.Index(password, ":")
			if index == -1 {
				goto invalidRequest
			}
			password = password[index+1:]
			if err == nil && password == config.Config.ControlPanel.Password {
				log.Println("[Control Panel] Successful login from " + req.RemoteAddr)
				cookie, err = m.sessions.FillRandom(util.AlphaNum, 32, nil)
				sendSessionCookie(res, cookie)
				handler(res, req)
				return
			}

			log.Println("[Control Panel] Bad login from " + req.RemoteAddr + " with wrong password '" + string(password) + "'")
			fn := m.failures.Get(host)
			if fn == nil {
				m.failures.Fill(host, func() interface{} { return new(failures) }, true)
				fn = m.failures.Get(host)
			}
			f := fn.(*failures)
			f.mutex.Lock()
			f.cnt++
			cnt := f.cnt
			f.mutex.Unlock()
			if cnt >= 10 {
				m.bans.Fill(host, func() interface{} { return nil }, true)
				log.Println("[Control Panel] " + host + " has been banned due to too many failed logins. Restart this process to reset.")
			}
		}

	invalidRequest:
		res.Header().Set("WWW-Authenticate", `Basic realm="Control panel password"`)
		res.WriteHeader(401)
		res.Write([]byte("Please login"))
		return
	}
}

func sendSessionCookie(res http.ResponseWriter, cookie string) {
	http.SetCookie(res, &http.Cookie{
		Name:     "RedoxCPSession",
		Path:     "/",
		Value:    cookie,
		HttpOnly: true,
		Secure:   config.Config.ControlPanel.SSL != nil,
		MaxAge:   int(config.Config.ControlPanel.SessionDuration.Seconds()),
	})
}
