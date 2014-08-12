package main

import (
	"testing"
)

func TestDecodeFile(t *testing.T) {

	decodeFile("test/Example.spritebuilder/SpriteBuilder Resources/MainScene.ccb")

}

func TestDecodeFileJSON(t *testing.T) {

	decodeFileJSON("test/Example.spritebuilder/SpriteBuilder Resources/MainScene.ccb")

}
