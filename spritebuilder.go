package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"howett.net/plist"
)

const (
	// memberName
	CONSTRUCTOR_TEMPLATE = "%s _%s(nullptr)\n"
	// memberName
	DESTRUCTOR_TEMPLATE = "CC_SAFE_RELEASE_NULL(_%s);\n"
	// MappingClassName memberName
	MEMBER_TEMPLATE = "%s* _%s;\n"
	// memberName MappingClassName memberName
	ASSIGN_CCB_MEMBER_TEMPLATE = "SB_MEMBERVARIABLEASSIGNER_GLUE(this, \"%s\", %s*, this->_%s);\n"
)

var CCBConvertClassMapping = map[string]string{
	"CCSprite9Slice": "cocos2d::extension::Scale9Sprite",
	"CCSprite":       "cocos2d::Sprite",
	"CCLabelTTF":     "cocos2d::Label",
}

type ccbRoot struct {
	UUID      int       `json:"UUID"`
	NodeGraph children  `json:"nodeGraph"`
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

func (c *children) Cocos2dxClassName() string {
	return CCBConvertClassMapping[c.BaseClass]
}

type properties struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type sequences struct {
	AutoPlay          bool `json:"autoPlay"`
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

func CreateCppCodeToCCBFile(filePath string) error {
	ccb, err := readCCBFile(filePath)
	if err != nil {
		return err
	}

	fmt.Println(strings.Repeat("-", 163))
	createDestructorCppCodeChildren(0, ccb.NodeGraph.Children)

	fmt.Println(strings.Repeat("-", 163))
	createCppMemberCodeChildren(0, ccb.NodeGraph.Children)

	fmt.Println(strings.Repeat("-", 163))
	createCppAssignCCBMemberCodeChildren(0, ccb.NodeGraph.Children)

	fmt.Println(strings.Repeat("-", 163))
	createConstructorCppCodeChildren(0, ccb.NodeGraph.Children)

	return nil
}

func createConstructorCppCodeChildren(count int, c []children) {
	for _, child := range c {
		if !notCustomNode || child.CustomClass != "" || child.MemberVarAssignmentName != "" {
			fmt.Printf(CONSTRUCTOR_TEMPLATE, ",", child.MemberVarAssignmentName)
		}
		createConstructorCppCodeChildren(count+1, child.Children)
	}
}

func createDestructorCppCodeChildren(count int, c []children) {
	for _, child := range c {
		if !notCustomNode || child.CustomClass != "" || child.MemberVarAssignmentName != "" {
			fmt.Printf(DESTRUCTOR_TEMPLATE, child.MemberVarAssignmentName)
		}
		createDestructorCppCodeChildren(count+1, child.Children)
	}
}

func createCppMemberCodeChildren(count int, c []children) {
	for _, child := range c {
		if !notCustomNode || child.CustomClass != "" || child.MemberVarAssignmentName != "" {
			fmt.Printf(MEMBER_TEMPLATE, child.Cocos2dxClassName(), child.MemberVarAssignmentName)
		}
		createCppMemberCodeChildren(count+1, child.Children)
	}
}

func createCppAssignCCBMemberCodeChildren(count int, c []children) {
	for _, child := range c {
		if !notCustomNode || child.CustomClass != "" || child.MemberVarAssignmentName != "" {
			fmt.Printf(ASSIGN_CCB_MEMBER_TEMPLATE, child.MemberVarAssignmentName, child.Cocos2dxClassName(), child.MemberVarAssignmentName)
		}
		createCppAssignCCBMemberCodeChildren(count+1, child.Children)
	}
}

func CheckReadCCBFile(filePath string) error {
	ccb, err := readCCBFile(filePath)
	if err != nil {
		return err
	}
	index := strings.LastIndex(filePath, "/")
	fmt.Printf("- %s\n", filePath[index+1:])
	fmt.Println(strings.Repeat("-", 163))
	fmt.Printf("| %-30s | %-40s | %-40s | %-40s |\n", "BaseClassName", "DisplyName", "CustomeClass", "MemberName")
	fmt.Println(strings.Repeat("-", 163))

	checkChildren(0, ccb.NodeGraph.Children)

	// timeline
	fmt.Println(strings.Repeat("-", 163))
	fmt.Printf("| %-30s | %-10s | %-40s | %-70s |\n", "TimelineName", "AutoPlay", "", "")
	fmt.Println(strings.Repeat("-", 163))
	for _, seq := range ccb.Sequences {
		autoPlay := ""
		if seq.AutoPlay {
			autoPlay = "ON"
		}
		fmt.Printf("| %-30s | %-10s | %-40s | %-70s |\n", seq.Name,  autoPlay, "", "")
	}

	if err := CreateCppCodeToCCBFile(filePath); err != nil {
		return err
	}

	return nil
}

func checkChildren(count int, c []children) {
	for _, child := range c {
		if !notCustomNode || child.CustomClass != "" || child.MemberVarAssignmentName != "" {
			baseName := strings.Join([]string{strings.Repeat(childCounterChar, count), child.BaseClass}, " ")
			fmt.Printf("| %-30s | %-40s | %-40s | %-40s |\n", baseName, child.DisplayName, child.CustomClass, child.MemberVarAssignmentName)
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

//	fmt.Println(string(j))

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
