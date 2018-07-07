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
	"sync"
	"time"
)

type ExpiringSyncMap struct {
	duration Duration
	mutex    sync.RWMutex
	values   map[string]entry
}

func NewExpiringSyncMap(expiry Duration) *ExpiringSyncMap {
	return &ExpiringSyncMap{
		duration: expiry,
		values:   make(map[string]entry),
	}
}

type entry struct {
	expiry time.Time
	value  interface{}
}

func (m *ExpiringSyncMap) syncCleanup() {
	for key, value := range m.values {
		if value.expiry.After(time.Now()) {
			delete(m.values, key)
		}
	}
}

func (m *ExpiringSyncMap) Overwrite(key string, value interface{}, cleanup bool) (overwritten bool) {
	m.mutex.Lock()
	if cleanup {
		m.syncCleanup()
	}
	_, overwritten = m.values[key]
	m.values[key] = entry{
		expiry: time.Now().Add(m.duration.Duration),
		value:  value,
	}
	m.mutex.Unlock()
	return
}

func (m *ExpiringSyncMap) Fill(key string, valueCreator func() interface{}, cleanup bool) (filled bool) {
	m.mutex.RLock()
	ntr, exists := m.values[key]
	m.mutex.RUnlock()
	if exists && ntr.expiry.After(time.Now()) {
		return false
	}

	value := valueCreator()
	expiry := time.Now().Add(m.duration.Duration)
	ntr = entry{
		expiry: expiry,
		value:  value,
	}
	m.mutex.Lock()
	if cleanup {
		m.syncCleanup()
	}
	_, exists = m.values[key]
	if !exists {
		m.values[key] = ntr
	}
	m.mutex.Unlock()
	return !exists
}

func (m *ExpiringSyncMap) FillRandom(charset []rune, size int, value interface{}) (key string, err error) {
	expiry := time.Now().Add(m.duration.Duration)
	ntr := entry{
		expiry: expiry,
		value:  value,
	}

	for {
		key, err = CryptoSecureRandomString(charset, size)
		if err != nil {
			return
		}

		m.mutex.Lock()
		_, exists := m.values[key]
		if exists {
			m.mutex.Unlock()
			continue
		}

		m.values[key] = ntr
		m.mutex.Unlock()
		return
	}
}

func (m *ExpiringSyncMap) Get(key string) (value interface{}) {
	m.mutex.RLock()
	ntr, exists := m.values[key]
	m.mutex.RUnlock()
	if exists && ntr.expiry.After(time.Now()) {
		return ntr.value
	} else {
		return nil
	}
}
func (m *ExpiringSyncMap) GetExists(key string) (value interface{}, exists bool) {
	m.mutex.RLock()
	ntr, exists := m.values[key]
	m.mutex.RUnlock()
	if exists && ntr.expiry.After(time.Now()) {
		return ntr.value, true
	} else {
		return nil, false
	}
}
