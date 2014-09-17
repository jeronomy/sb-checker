package main

import (
	"testing"
)

func TestCheckReadCCBFile(t *testing.T) {
	testCCBFilePath := "test/Example.spritebuilder/SpriteBuilder Resources/MainScene.ccb"
	if err := CheckReadCCBFile(testCCBFilePath); err != nil {
		t.Error(err)
	}
}
func TestCheckReadCCBDir(t *testing.T) {
	testCCBDirPath := "test/Example.spritebuilder/SpriteBuilder Resources/"
	if err := CheckReadCCBDir(testCCBDirPath); err != nil {
		t.Error(err)
	}
}
