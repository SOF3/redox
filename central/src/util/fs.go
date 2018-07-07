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
	"os"
	"path"
)

func FileExists(file string) (exists bool, err error) {
	info, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	exists = !info.IsDir()
	return
}
func DirExists(file string) (exists bool, err error) {
	info, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	exists = info.IsDir()
	return
}

func SafeMkdir(dir string, recursive bool) (err error) {
	var exists bool
	if exists, err = DirExists(dir); err != nil || exists {
		return
	}

	if recursive {
		err = os.MkdirAll(dir, 0755)
	} else {
		err = os.Mkdir(dir, 0755)
	}

	return
}

func SafeOpenWrite(file string) (*os.File, error) {
	dir := path.Dir(file)
	if err := SafeMkdir(dir, true); err != nil {
		return nil, err
	}

	return os.Create(file)
}

func MustOpenWrite(file string) *os.File {
	f, err := SafeOpenWrite(file)
	if err != nil {
		panic(err)
	}
	return f
}
