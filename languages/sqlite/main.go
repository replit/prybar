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

	return f.Name()
}

// ripped from goval
// Set some modes on the pty attached to fd such that it is
// more suited to be connected to a terminal.
func SaneTTY(fd int) (interface{}, error) {
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}

	termios.Iflag |= unix.ICRNL
	termios.Oflag |= unix.OPOST
	termios.Lflag |= unix.ECHO | unix.ISIG | unix.IEXTEN | unix.ICANON
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, termios); err != nil {
		return nil, err
	}
	return termios, nil
}

// ripped from goval
// Disable echo on the pty attached to fd
func EchoOffTTY(fd int) (interface{}, error) {
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}

	termios.Lflag &^= unix.ECHO
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, termios); err != nil {
		return nil, err
	}
	return termios, nil
}

func Execute(config *utils.Config) {
	configFile := constructConfigFile(config)
	args := []string{"-init", configFile}
	cmd := exec.Command("sqlite", args...)

	ptty, tty, err := pty.Open()
	if err != nil {
		panic(err)
	}

	SaneTTY(int(ptty.Fd()))
	EchoOffTTY(int(ptty.Fd()))
	cmd.Stderr = tty
	cmd.Stdin = tty
	cmd.Stdout = tty

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	// file to execute
	if len(config.Args) > 1 {
		panic("too many arguments")
	}
	if len(config.Args) == 1 {
		fileToRun := config.Args[0]
		ptty.WriteString(fmt.Sprintf(".read %s\n", fileToRun))
	}
	
	// set up I/O
	go io.Copy(os.Stderr, ptty)
	if config.Quiet {
		go io.Copy(os.Stdout, filter(ptty))
	} else {
		go io.Copy(os.Stdout, ptty)
	}
	io.Copy(ptty, os.Stdin)
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
