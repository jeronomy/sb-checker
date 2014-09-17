package main

import (
	"fmt"
	"io/ioutil"
	"strings"
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
	"CCButton":       "cocos2d::extension::ControlButton",
	"CCBFile":        "cocos2d::Node",
	"CCNode":         "cocos2d::Node",
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
		fmt.Printf("| %-30s | %-10s | %-40s | %-70s |\n", seq.Name, autoPlay, "", "")
	}

	if err := CreateCppCodeToCCBFile(filePath); err != nil {
		return err
	}

	return nil
}

func checkChildren(count int, c []children) {
	for _, child := range c {
		if !notCustomNode || child.CustomClass != "" || child.MemberVarAssignmentName != "" {
			customClass := child.getCCBCustomClass()
			baseName := strings.Join([]string{strings.Repeat(childCounterChar, count), child.BaseClass}, " ")
			fmt.Printf("| %-30s | %-40s | %-40s | %-40s |\n", baseName, child.DisplayName, customClass, child.MemberVarAssignmentName)
		}
		checkChildren(count+1, child.Children)
	}
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
			fmt.Printf(MEMBER_TEMPLATE, child.getCocos2dxClassName(), child.MemberVarAssignmentName)
		}
		createCppMemberCodeChildren(count+1, child.Children)
	}
}

func createCppAssignCCBMemberCodeChildren(count int, c []children) {
	for _, child := range c {
		if !notCustomNode || child.CustomClass != "" || child.MemberVarAssignmentName != "" {
			fmt.Printf(ASSIGN_CCB_MEMBER_TEMPLATE, child.MemberVarAssignmentName, child.getCocos2dxClassName(), child.MemberVarAssignmentName)
		}
		createCppAssignCCBMemberCodeChildren(count+1, child.Children)
	}
}
