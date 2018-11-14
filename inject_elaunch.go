// this file will be injected into a language's main package at build.
package main

import "github.com/replit/prybar/utils"

func main() {
	Execute(utils.ParseFlags())
}
