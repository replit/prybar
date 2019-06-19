package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/replit/prybar/utils"
)

// constructConfigFile generates commands to configure the sqlite CLI.
// It writes them to a temporary file and returns its pathname.
func constructConfigFile(config *utils.Config) string {
	f, err := ioutil.TempFile("", "sqlite-config")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// main and continuation prompts
	// TODO: this probably doesn't handle quotation marks properly
	f.WriteString(fmt.Sprintf(".prompt '%s' '%s'\n", config.Ps1, config.Ps2))

	return f.Name()
}

func Execute(config *utils.Config) {
	sqlite, err := exec.LookPath("sqlite")
	if err != nil {
		panic(err)
	}

	configFile := constructConfigFile(config)

	args := []string{"sqlite", "-init", configFile}
	err = syscall.Exec(sqlite, args, os.Environ())
	if err != nil {
		panic(err)
	}
}
