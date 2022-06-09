package main

import (
	"testing"
)

// Mostly a debugging aid, not sure if it'll actually fail if something goes wrong
func TestLuaStuff(t *testing.T) {
	var lua Lua
	lua.Open()
	lua.Version()
	lua.EvalFile("testdata/test.lua", make([]string, 0))
}
