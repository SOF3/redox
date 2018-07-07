<?php

/*
 * redox/pocketmine
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

declare(strict_types=1);

namespace Redox\PocketMine;

use RuntimeException;
use function chr;
use function fwrite;
use function stream_set_blocking;
use function stream_socket_client;
use function strlen;
use const PHP_INT_SIZE;

class RedoxConnector{
	/** @var array */
	protected $config;
	/** @var resource */
	protected $socket;

	public function __construct(array $config){
		$this->config = $config;
		$this->socket = stream_socket_client("ssl://" . $this->config["host"] . ":" . $this->config["port"], $errNo, $errStr);
		stream_set_blocking($this->socket, 1);

		$this->writeLInt(2, RedoxProtocol::ID_CLIENT_HANDSHAKE);
		$this->writeLInt(2, RedoxProtocol::VERSION);
		$this->writeString($this->config["node-id"]);
		$this->writeString($this->config["password"]);

		$this->readPacket();
	}

	protected function writeLInt(int $bytes, int $number) : void{
		// little-endian byte order
		if($bytes > PHP_INT_SIZE){
			throw new RuntimeException("cannot express such a large number in PHP");
		}
		$buffer = "";
		for($i = 0; $i < $bytes; ++$i){
			$ord = $number & 0xFF;
			$number >>= 8;
			$buffer .= chr($ord);
		}

		fwrite($this->socket, $buffer);
	}

	public function writeString(string $string) : void{
		$this->writeLInt(2, strlen($string));
		fwrite($this->socket, $string);
	}
}
