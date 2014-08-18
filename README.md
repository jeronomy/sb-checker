sb-checkeer

[![Build Status](https://drone.io/github.com/kyokomi/sb-checker/status.png)](https://drone.io/github.com/kyokomi/sb-checker/latest) 
[![Coverage Status](https://img.shields.io/coveralls/kyokomi/sb-checker.svg)](https://coveralls.io/r/kyokomi/sb-checker?branch=master)


Integrity check tool of [SpriteBuilder](http://www.spritebuilder.com/) for golang（Go）

====

## Usage

```sh
$ sb-checker help
NAME:
   sb-checkeer -

USAGE:
   sb-checkeer [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
  kyokomi - <kyoko1220adword@gmail.com>

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input-ccb 		input spritebuilder ccb file path [$INPUT_CCB_FILE_PATH]
   --input-ccb-dir 	input spritebuilder ccb directry path [$INPUT_CCB_DIR_PATH]
   -d
   --help, -h		show help
   --version, -v	print the version
```

## Demo

```sh
$ sb-checker --input-ccb-dir test/Example.spritebuilder/SpriteBuilder\ Resources/
-MainScene.ccb
|  CCNodeGradient                | FugaNode                                 | fugaNode                                 |
|  CCLabelTTF
```

## Install

comming soon...

## Licence

[MIT](https://github.com/kyokomi/sb-checker/blob/master/LICENSE)

## Author

[kyokomi](https://github.com/kyokomi)

