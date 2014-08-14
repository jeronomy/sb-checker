// package main
package main

import (
	"os"
	"strings"

	"log"

	"encoding/json"

	"fmt"

	"github.com/codegangsta/cli"
	"howett.net/plist"
)

type ccbRoot struct {
	UUID      int       `json:"UUID"`
	NodeGraph nodeGraph `json:"nodeGraph"`
}
type nodeGraph struct {
	UUID      int        `json:"UUID"`
	BaseClass string     `json:"baseClass"`
	Childrens []children `json:"children"`
}
type children struct {
	UUID                    int         `json:"UUID"`
	BaseClass               string      `json:"baseClass"`
	Childrens               []children  `json:"children"`
	CustomClass             string      `json:"customClass"`
	DisplayName             string      `json:"displayName"`
	MemberVarAssignmentName string      `json:"memberVarAssignmentName"`
	MemberVarAssignmentType int         `json:"memberVarAssignmentType"`
	Properties              []propertie `json:"properties"`
}
type propertie struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
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

const (
	childCounterChar = "->"
)

// $ go run sb-checker.go version.go --input-ccb="/Users/kyokomi/src/github.com/kyokomi/sb-checkeer/test/Example.spritebuilder/SpriteBuilder Resources/MainScene.ccb"
func doMain(c *cli.Context) {

	ccb, err := readCCBFile(c.String("input-ccb"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ccb)
	fmt.Println("BaseClass = ", ccb.NodeGraph.BaseClass)
	checkChildrens(0, ccb.NodeGraph.Childrens)
}

func checkChildrens(count int, childrens []children) {
	for _, child := range childrens {
		fmt.Println(strings.Repeat(childCounterChar, count), "customClass = ", child.CustomClass)
		for _, prop := range child.Properties {
			if prop.Name == "name" {
				fmt.Println(strings.Repeat(childCounterChar, count), "name = ", prop.Value.(string))
			}
		}
		checkChildrens(count+1, child.Childrens)
	}
}

// ccbファイルを読み込む。
func readCCBFile(filePath string) (*ccbRoot, error) {
	j, err := decodeFileJSON(filePath)
	if err != nil {
		return nil, err
	}

	var ccb ccbRoot
	if err := json.Unmarshal(j, &ccb); err != nil {
		return nil, err
	}
	return &ccb, nil
}

// plistのファイルをjsonに変換する.
func decodeFileJSON(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := plist.NewDecoder(f)
	var m = make(decodeCcbRoot)
	if err := d.Decode(m); err != nil {
		return nil, err
	}

	j, err := json.Marshal(&m)
	if err != nil {
		return nil, err
	}
	return j, nil
}
