package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/kr/pty"
	"golang.org/x/sys/unix"

	"github.com/replit/prybar/utils"
)

var Instance = &SQLite{}

type SQLite struct{}

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

func disableEcho(file *os.File) {
	fd := int(file.Fd())
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		panic(err)
	}

	termios.Lflag &^= unix.ECHO

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, termios); err != nil {
		panic(err)
	}
}

func Execute(config *utils.Config) {
	// sanity check
	if len(config.Args) > 1 {
		panic("too many arguments")
	}

	configFile := constructConfigFile(config)
	args := []string{"-init", configFile}
	cmd := exec.Command("sqlite3", args...)

	ptty, tty, err := pty.Open()
	if err != nil {
		panic(err)
	}
	disableEcho(ptty)

	// don't hook up stderr until after our config is loaded to avoid unnecessary output
	cmd.Stderr = ioutil.Discard
	cmd.Stdin = os.Stdin
	cmd.Stdout = tty

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	if config.Quiet {
		waitForPrompt(ptty, config.Ps1)
	}

	// now that we have a prompt, we can hook up stderr
	cmd.Stderr = os.Stderr

	// execute file, if specified
	if len(config.Args) == 1 {
		fileToRun := config.Args[0]
		ptty.WriteString(fmt.Sprintf(".read %s\n", fileToRun))
	}

	io.Copy(os.Stdout, ptty)
}

func waitForPrompt(src io.Reader, prompt string) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		if scanner.Text() == "Enter \".help\" for usage hints." {
			return
		}
	}
}
