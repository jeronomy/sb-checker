package main

import "testing"

const (
	testCCBFilePath = "test/Example.spritebuilder/SpriteBuilder Resources/MainScene.ccb"
)

func TestReadCCBFile(t *testing.T) {
	ccb, err := readCCBFile(testCCBFilePath)
	if err != nil {
		t.Error(err)
	}
	if ccb.NodeGraph.BaseClass != "CCNode" {
		t.Error("baseClass is not CCNode > ", ccb.NodeGraph.BaseClass)
	}
}

func TestDecodeFileJSON(t *testing.T) {
	if _, err := decodeFileJSON(testCCBFilePath); err != nil {
		t.Error(err)
	}
}
