package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"howett.net/plist"
)

type ccbRoot struct {
	UUID      int       `json:"UUID"`
	NodeGraph nodeGraph `json:"nodeGraph"`
}

type nodeGraph struct {
	UUID      int        `json:"UUID"`
	BaseClass string     `json:"baseClass"`
	Children  []children `json:"children"`
}

type children struct {
	UUID                    int          `json:"UUID"`
	BaseClass               string       `json:"baseClass"`
	Children                []children   `json:"children"`
	CustomClass             string       `json:"customClass"`
	DisplayName             string       `json:"displayName"`
	MemberVarAssignmentName string       `json:"memberVarAssignmentName"`
	MemberVarAssignmentType int          `json:"memberVarAssignmentType"`
	Properties              []properties `json:"properties"`
}

type properties struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type decodeCcbRoot map[string]interface{}

var notCustomNode = false

const (
	childCounterChar = "->"
)

func CheckReadCCBDir(dirPath string) error {
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range fs {
		filePath := strings.Join([]string{dirPath, file.Name()}, "/")
		if file.IsDir() {
			if err := CheckReadCCBDir(filePath); err != nil {
				return err
			}
		} else if strings.HasSuffix(file.Name(), ".ccb") {
			if err := CheckReadCCBFile(filePath); err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckReadCCBFile(filePath string) error {
	ccb, err := readCCBFile(filePath)
	if err != nil {
		return err
	}
	index := strings.LastIndex(filePath, "/")
	fmt.Printf("- %s\n", filePath[index+1:])
	checkChildren(0, ccb.NodeGraph.Children)

	return nil
}

func checkChildren(count int, c []children) {
	for _, child := range c {
		if !notCustomNode || child.CustomClass != "" || child.MemberVarAssignmentName != "" {
			baseName := strings.Join([]string{strings.Repeat(childCounterChar, count), child.BaseClass}, " ")
			fmt.Printf("| %-30s | %-40s | %-40s |\n", baseName, child.CustomClass, child.MemberVarAssignmentName)
		}
		checkChildren(count+1, child.Children)
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
