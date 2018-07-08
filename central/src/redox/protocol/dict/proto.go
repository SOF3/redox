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
	"encoding/binary"
	"io"
)

const (
	UnitHeader      byte = 0x00
	UnitSetDictSpec byte = 0x01
	UnitDefineWord  byte = 0x02
	UnitMessage     byte = 0x03
)

var HeaderMagic = []byte{
	0xff,
	0x72,
	0x65,
	0x64,
	0x6f,
	0x78,
	0x00,
	0x64,
	0x69,
	0x63,
	0x74,
	0xfe,
}

const HeaderVersion uint16 = 1

type iOutput interface {
	io.ByteWriter
	io.Writer
}

type Sender struct {
	Output iOutput
	Dict   *Dictionary
}

func New(output iOutput) (sender *Sender, err error) {
	sender = &Sender{
		Output: output,
		Dict:   NewDictionary(),
	}

	if err = sender.sendUnitHeader(); err != nil {
		return
	}

	return
}

func (sender *Sender) SetSpec(spec Spec) (err error) {
	if err = sender.SetSpec(spec); err != nil {
		return
	}
	if err = sender.sendUnitSetDictSpec(spec); err != nil {
		return
	}
	return
}

func (sender *Sender) SendMessage(messageType string, messageData []byte) (err error) {
	id, send, err := sender.Dict.UseWord(messageType)
	if err != nil {
		return
	}

	if send {
		if err = sender.sendUnitDefineWord(messageType); err != nil {
			return
		}
	}

	if err = sender.sendUnitMessage(id, messageData); err != nil {
		return
	}

	return
}

func (sender *Sender) sendUvarint(i uint64) (err error) {
	buf := make([]byte, binary.MaxVarintLen64)
	length := binary.PutUvarint(buf, i)
	_, err = sender.Output.Write(buf[:length])
	return
}

func (sender *Sender) sendUnitHeader() (err error) {
	if err = sender.Output.WriteByte(UnitHeader); err != nil {
		return
	}
	if _, err = sender.Output.Write(HeaderMagic); err != nil {
		return
	}
	if err = binary.Write(sender.Output, binary.BigEndian, HeaderVersion); err != nil {
		return
	}
	return
}

func (sender *Sender) sendUnitSetDictSpec(spec Spec) (err error) {
	if err = sender.Output.WriteByte(UnitSetDictSpec); err != nil {
		return
	}
	if err = sender.sendUvarint(uint64(spec.TraceSize)); err != nil {
		return
	}
	if err = sender.sendUvarint(uint64(spec.Frequency)); err != nil {
		return
	}
	return
}

func (sender *Sender) sendUnitDefineWord(word string) (err error) {
	if err = sender.Output.WriteByte(UnitDefineWord); err != nil {
		return
	}
	if err = sender.sendUvarint(uint64(len(word))); err != nil {
		return
	}
	if _, err = sender.Output.Write([]byte(word)); err != nil {
		return
	}
	return
}

func (sender *Sender) sendUnitMessage(typeId WordId, messageData []byte) (err error) {
	if err = sender.Output.WriteByte(UnitMessage); err != nil {
		return
	}
	if err = sender.sendUvarint(uint64(typeId)); err != nil {
		return
	}
	if _, err = sender.Output.Write(messageData); err != nil {
		return
	}
	return
}

type Recipient struct {
	Input io.ByteReader
	Dict  Dictionary
}
