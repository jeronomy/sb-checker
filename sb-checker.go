// package main
package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

// $ sb-checker --input-ccb="test/Example.spritebuilder/SpriteBuilder Resources/MainScene.ccb"
func doMain(c *cli.Context) {

	notCustomNode = c.Bool("d")

	if c.String("input-ccb-dir") == "" {
		if err := CheckReadCCBFile(c.String("input-ccb")); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := CheckReadCCBDir(c.String("input-ccb")); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "sb-checkeer"
	app.Version = Version
	app.Usage = ""
	app.Author = "kyokomi"
	app.Email = "kyoko1220adword@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{"input-ccb", "", "input spritebuilder ccb file path", "INPUT_CCB_FILE_PATH"},
		cli.StringFlag{"input-ccb-dir", "", "input spritebuilder ccb directry path", "INPUT_CCB_DIR_PATH"},
		cli.BoolFlag{"d", "", ""},
	}
	app.Action = doMain
	app.Run(os.Args)
}
