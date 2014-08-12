package main

import (
	"os"

	"fmt"
	"log"

	"encoding/json"

	"github.com/codegangsta/cli"
	"howett.net/plist"
)

type ccbRoot struct {
	uuid      int       `json:UUID`
	nodeGraph nodeGraph `json:"nodeGraph"`
}
type nodeGraph struct {
	uuid      int        `json:UUID`
	childrens []children `json:"children"`
}
type children struct {
	customeName string      `json:"customeName"`
	properties  []propertie `json:"propertie"`
}
type propertie struct {
	name      string      `json:"name"`
	value     interface{} `json:"value"`
	childrens []children  `json:"children"`
}

type decodeCcbRoot map[string]interface{}

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

// $ go run sb-checkeer.go version.go --input-ccb="/Users/kyokomi/src/github.com/kyokomi/sb-checkeer/test/Example.spritebuilder/SpriteBuilder Resources/MainScene.ccb"
func doMain(c *cli.Context) {

	//	decodeFile(c.String("input-ccb"))
	decodeFileJSON(c.String("input-ccb"))
}

func decodeFileJSON(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := plist.NewDecoder(f)
	var m = make(decodeCcbRoot)
	if err := d.Decode(m); err != nil {
		log.Fatal(err)
	}

	j, err := json.Marshal(&m)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(j))
}

func decodeFile(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := plist.NewDecoder(f)
	var m = make(map[string]interface{})
	if err := d.Decode(m); err != nil {
		log.Fatal(err)
	}

	var nodes = make(map[string]interface{})
	nodes = m["nodeGraph"].(map[string]interface{})
	readChildren(nodes["children"].([]interface{}))
}

func readChildren(children []interface{}) {
	for _, child := range children {
		child := child.(map[string]interface{})

		fmt.Println("customeName:", child["customClass"])

		for _, properties := range child["properties"].([]interface{}) {
			p := properties.(map[string]interface{})
			if p["name"].(string) == "name" {
				fmt.Println("name:", p["value"])
			}
		}

		readChildren(child["children"].([]interface{}))
	}
}
