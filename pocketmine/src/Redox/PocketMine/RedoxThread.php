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

use function chr;
use pocketmine\Thread;
use RuntimeException;
use function socket_connect;
use function socket_create;
use const AF_INET;
use const PHP_INT_SIZE;
use const SOCK_STREAM;
use function socket_write;
use const SOL_TCP;

class RedoxThread extends Thread{
	/** @var array */
	protected $config;
	/** @var bool */
	protected $stopped = false;

	public function __construct(array $config){
		$this->config = $config;
	}

	public function close() : void{
		$this->stopped = true;
		$this->join();
	}

	public function run() : void{
		$conn = new RedoxConnector($this->config);

	}
}
