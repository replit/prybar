package main

// USING_CGO

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"

	"github.com/replit/prybar/utils"
)

var Instance = SQLite{}

// some compile-time assertions that we satisfy interfaces
var _ utils.PluginBase = SQLite{}
var _ utils.PluginEval = SQLite{}
var _ utils.PluginEvalExpression = SQLite{}


type SQLite struct {
	open bool
	db *sql.DB
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

func (s SQLite) Version() string {
	return "SQLite version TODO"
}

func (s SQLite) Open() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	s.db = db
	s.open = true
}
func (s SQLite) Close() {
	s.db.Close()
}

func (s SQLite) Eval(line string) {
	if !s.open {
		s.Open()
	}
	_, err := s.db.Exec(line)
	if err != nil {
		panic(err)
	}
}

func (s SQLite) EvalExpression(line string) string {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(line)
	// rows, err := s.db.Query(line)
	if err != nil {
		return err.Error()
	}
	return "evaled"
}
