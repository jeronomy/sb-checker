// package main
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"howett.net/plist"
	"io/ioutil"
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

var notCustomNode = false

func main() {
	app := cli.NewApp()
	app.Name = "sb-checkeer"
	app.Version = Version
	app.Usage = ""
	app.Author = "kyokomi"
	app.Email = "kyoko1220adword@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{"input-ccb",     "", "input spritebuilder ccb file path", "INPUT_CCB_FILE_PATH"},
		cli.StringFlag{"input-ccb-dir", "", "input spritebuilder ccb directry path", "INPUT_CCB_DIR_PATH"},
		cli.BoolFlag{"d", "", ""},
	}
	app.Action = doMain
	app.Run(os.Args)
}

const (
	childCounterChar = "->"
)

// $ go run sb-checker.go version.go --input-ccb="/Users/kyokomi/src/github.com/kyokomi/sb-checkeer/test/Example.spritebuilder/SpriteBuilder Resources/MainScene.ccb"
func doMain(c *cli.Context) {

	notCustomNode = c.Bool("d")

	if c.String("input-ccb-dir") == "" {
		ccb, err := readCCBFile(c.String("input-ccb"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("BaseClass = ", ccb.NodeGraph.BaseClass)
		checkChildrens(0, ccb.NodeGraph.Childrens)
	} else {

		err := checkReadCCBDir(c.String("input-ccb-dir"))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func checkReadCCBDir(dirPath string) error {
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range fs {
		if file.IsDir() {
			err := checkReadCCBDir(strings.Join([]string{dirPath, file.Name()} , "/"))
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(file.Name(), ".ccb") {
			ccb, err := readCCBFile(strings.Join([]string{dirPath, file.Name()} , "/"))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("-%s\n", file.Name())
			checkChildrens(0, ccb.NodeGraph.Childrens)
		}
	}
	return nil
}

func checkChildrens(count int, childrens []children) {
	for _, child := range childrens {
		if !notCustomNode || child.CustomClass != "" || child.MemberVarAssignmentName != "" {
			baseName := strings.Join([]string{strings.Repeat(childCounterChar, count), child.BaseClass}, " ")
			fmt.Printf("| %-30s | %-40s | %-40s |\n", baseName, child.CustomClass, child.MemberVarAssignmentName)
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
