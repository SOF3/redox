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

package dict

import (
	"container/list"
	"errors"
	"github.com/SOF3/redox/central/src/util"
	"sort"
)

type WordId uint

type Dictionary struct {
	specSet bool
	spec    Spec

	wordList []string
	wordMap  map[string]WordId
	trace    *list.List
	counter  int
}

type Spec struct {
	TraceSize int
	Frequency int
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		wordList: []string{},
		wordMap:  map[string]WordId{},
	}
}

func (dict *Dictionary) SetSpec(spec Spec) (err string) {
	if spec.TraceSize > 4096 {
		return "The other party sent an abnormally high TraceSize, possibly a DoS attempt"
	}

	dict.spec = spec
	dict.specSet = true
	dict.trace = list.New()
	dict.counter = 0
	return
}

func (dict *Dictionary) Define(word string) (id WordId, err string) {
	if !dict.specSet {
		return -1, "Spec not set"
	}

	if _, exists := dict.wordMap[word]; exists {
		return -1, "word " + word + " already defined"
	}
	id = WordId(len(dict.wordList))
	dict.wordList = append(dict.wordList, word)
	dict.wordMap[word] = id
	return
}

func (dict *Dictionary) Use(id WordId) (word string, err string) {
	if !dict.specSet {
		return "", "Spec not set"
	}

	dict.trace.PushBack(id)
	if dict.trace.Len() > dict.spec.TraceSize {
		dict.trace.Remove(dict.trace.Front())
	}
	util.Assert(dict.trace.Len() <= dict.spec.TraceSize, "Dictionary trace got too many values")

	dict.counter++
	util.Assert(dict.counter <= dict.spec.Frequency, "Dictionary counter exceeded frequency")

	word = dict.wordList[id]

	if dict.counter == dict.spec.Frequency {
		dict.counter = 0
		freq := map[WordId]uint{}
		for el := dict.trace.Front(); el != nil; el = el.Next() {
			freq[el.Value.(WordId)]++
		}
		sort.SliceStable(dict.wordList, func(i, j int) bool {
			return freq[WordId(i)] > freq[WordId(j)]
		})
	}

	return
}

func (dict *Dictionary) UseWord(word string) (id WordId, new bool, e error) {
	id, exists := dict.wordMap[word]
	var err string
	if !exists {
		id, err = dict.Define(word)
		if err != "" {
			return -1, false, errors.New(err)
		}
	}

	_, err = dict.Use(id)
	if err != "" {
		return -1, false, errors.New(err)
	}
	return id, !exists, nil
}
