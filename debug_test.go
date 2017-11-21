package main

import (
	"errors"
	"testing"
)

func TestDebug_InitDebug(t *testing.T) {
	InitDebug(true)
	DebugHandler.Sys("Output", "Test")
	DebugHandler.Err(errors.New("Output"), "Test", 3)
}
