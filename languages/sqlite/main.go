package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/kr/pty"

	"github.com/replit/prybar/utils"
)

var Instance = &SQLite{}

type SQLite struct {}

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

	// version header (this is actually completely incorrect)
	var headers string
	if config.Quiet {
		headers = "off"
	} else {
		headers = "on"
	}
	f.WriteString(fmt.Sprintf(".headers %s\n", headers))

	return f.Name()
}

func Execute(config *utils.Config) {
	configFile := constructConfigFile(config)
	args := []string{"-init", configFile}
	cmd := exec.Command("sqlite", args...)

	ppty, tty, err := pty.Open()
	if err != nil {
		panic(err)
	}
	cmd.Stderr = tty
	cmd.Stdin = tty
	cmd.Stdout = tty
	_ = ppty

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	
	go io.Copy(os.Stderr, ppty)
	go io.Copy(os.Stdout, filter(ppty))
	io.Copy(ppty, os.Stdin)
}

// filter removes all output until we get a prompt
func filter(src io.Reader) io.Reader {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		if scanner.Text() == "Enter \".help\" for instructions" {
			break
		}
	}
	return src
}
