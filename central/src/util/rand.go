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
	"crypto/rand"
	"math/big"
)

var Alpha = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz")
var Num = []rune("01234576789")
var AlphaNum = append(Alpha, Num...)

func CryptoSecureRandomString(charset []rune, size int) (str string, err error) {
	length := big.NewInt(int64(len(charset)))

	output := make([]rune, size)

	for i := 0; i < size; i++ {
		var num *big.Int
		num, err = rand.Int(rand.Reader, length)
		if err != nil {
			return
		}

		output[i] = charset[num.Int64()]
	}

	str = string(output)
	return
}
