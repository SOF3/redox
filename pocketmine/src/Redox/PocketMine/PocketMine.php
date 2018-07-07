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

use pocketmine\plugin\PluginBase;
use function file_get_contents;
use function file_put_contents;
use function is_dir;
use function is_file;
use function mkdir;
use function yaml_emit;
use function yaml_parse;

class PocketMine extends PluginBase{
	/** @var array */
	protected $cfg;
	/** @var RedoxThread */
	protected $conn;

	public function onEnable() : void{
		if(!is_dir($this->getDataFolder())){
			mkdir($this->getDataFolder());
		}

		if(!is_file($this->getDataFolder() . "config.yml")){
			$this->getLogger()->critical("Please setup Redox by editing " . $this->getDataFolder() . "config.yml and start the server again");
			file_put_contents($this->getDataFolder() . "config.yml", yaml_emit([
				"README" => "Please refer to https://sof3.github.io/redox/pocketmine/config for instructions to edit this file.",
				"host" => "127.0.0.1",
				"port" => 45811,
				"password" => "",
				"node-id" => $this->getServer()->getServerUniqueId(),
				"shutdown-if-unconnected" => true,
				"Finally" => "Delete this line if you have finished editing this file.",
			]));
			$this->getServer()->shutdown();
			return;
		}

		$this->cfg = yaml_parse(file_get_contents($this->getDataFolder() . "config.yml"));
		if(isset($this->cfg["Finally"])){
			$this->getLogger()->critical("Please setup Redox by editing " . $this->getDataFolder() . "config.yml and start the server again. Remember to delete the \"Finally\" line.");
			return;
		}

		$this->conn = new RedoxThread($this->cfg);
		$this->conn->start();
	}
}
