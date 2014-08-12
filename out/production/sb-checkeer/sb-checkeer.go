package main

import (
	"os"

	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"howett.net/plist"
)

func main() {
	app := cli.NewApp()
	app.Name = "sb-checkeer"
	app.Version = Version
	app.Usage = ""
	app.Author = "kyokomi"
	app.Email = "kyoko1220adword@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{"input-ccb", "", "input spritebuilder ccb file path", "INPUT_CCB_FILE_PATH"},
	}
	app.Action = doMain
	app.Run(os.Args)
}

func doMain(c *cli.Context) {

	f, err := os.Open(c.String("input-ccb"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := plist.NewDecoder(f)
	var m = make(map[string]interface{})
	if err := d.Decode(m); err != nil {
		log.Fatal(err)
	}

	fmt.Println(tree)
	for _, t := range tree["nodeGraph"].(plist.Dict) {
		fmt.Println(t)
		if item, ok := t.(plist.Dict); ok {
			fmt.Println(item)
			fmt.Println(item["customClass"])
		}
	}
}
