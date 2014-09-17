package main

import (
	"encoding/json"
	"os"
	"howett.net/plist"
)

type ccbRoot struct {
	UUID      int         `json:"UUID"`
	NodeGraph children    `json:"nodeGraph"`
	Sequences []sequences `json:"sequences"`
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

type sequences struct {
	AutoPlay        bool `json:"autoPlay"`
	CallbackChannel struct {
		Keyframes []interface{} `json:"keyframes"`
		Type      float64       `json:"type"`
	} `json:"callbackChannel"`
	ChainedSequenceId float64 `json:"chainedSequenceId"`
	Length            float64 `json:"length"`
	Name              string  `json:"name"`
	Offset            float64 `json:"offset"`
	Position          float64 `json:"position"`
	Resolution        float64 `json:"resolution"`
	Scale             float64 `json:"scale"`
	SequenceId        float64 `json:"sequenceId"`
	SoundChannel      struct {
		IsExpanded bool          `json:"isExpanded"`
		Keyframes  []interface{} `json:"keyframes"`
		Type       float64       `json:"type"`
	} `json:"soundChannel"`
}

func (c *children) getCocos2dxClassName() string {
	customClass := c.getCCBCustomClass()
	if customClass != "" {
		return customClass
	}
	return CCBConvertClassMapping[c.BaseClass]
}

func (c *children) getCCBFileName() string {
	for _, properties := range c.Properties {
		if properties.Type != "CCBFile" {
			continue
		}
		return properties.Value.(string)
	}
	return ""
}

func (c *children) getCCBCustomClass() string {
	if c.CustomClass != "" {
		return c.CustomClass
	}

	filePath := c.getCCBFileName()
	if filePath == "" {
		return ""
	}

	// TODO: ccbRootDirectory currentDir
	ccb, err := readCCBFile(filePath)
	if err != nil {
		return ""
	}

	return ccb.NodeGraph.CustomClass
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
