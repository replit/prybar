package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/replit/prybar/utils"
)

type outputLine struct {
	line string
	source string
}

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

	// version header
	var headers string
	if config.Quiet {
		headers = "off"
	} else {
		headers = "on"
	}
	f.WriteString(fmt.Sprintf(".headers %s\n", headers))

	return f.Name()
}

func filterOutput(stdout io.Reader, stderr io.Reader, config *utils.Config) {
	outputChan := make(chan *outputLine)
	// start reading from the pipes
	go func() {
		s := bufio.NewScanner(stdout)
		for s.Scan() {
			line := s.Text()
			ol := &outputLine{line: line, source: "stdout"}
			outputChan <- ol
		}
		if err := s.Err(); err != nil {
			panic(err)
		}
	}()
	go func() {
		s := bufio.NewScanner(stderr)
		for s.Scan() {
			line := s.Text()
			ol := &outputLine{line: line, source: "stderr"}
			outputChan <- ol
		}
		if err := s.Err(); err != nil {
			panic(err)
		}
	}()

	for {
		fmt.Printf("top of for")
		for ol := range outputChan {
			fmt.Fprintf(os.Stdout, "%+v", ol)
		}
		panic("output channel closed")
	}
}

func Execute(config *utils.Config) {
	sqlite, err := exec.LookPath("sqlite")
	if err != nil {
		panic(err)
	}

	// set up sqlite command
	configFile := constructConfigFile(config)
	cmd := exec.Command(sqlite, "-init", configFile)
	cmd.Stdin = os.Stdin

	// grab output pipes
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// run sqlite and filter its output
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	filterOutput(cmdOut, cmdErr, config)
}
