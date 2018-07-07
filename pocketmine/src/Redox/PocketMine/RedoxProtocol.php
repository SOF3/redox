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

class RedoxProtocol{
	public const VERSION = 0x0001;
	public const ID_PING = 0x0001;
	public const ID_PONG = 0x0002;
	public const ID_CLIENT_HANDSHAKE = 0x0003;
	public const ID_SERVER_HANDSHAKE = 0x0004;
	public const ID_SERVER_DISCONNECT = 0x0005;
}
