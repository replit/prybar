package main

// USING_CGO

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/replit/prybar/utils"
)

var Instance = &SQLite{}

// some compile-time assertions that we satisfy interfaces
var _ utils.PluginBase = &SQLite{}
var _ utils.PluginEval = &SQLite{}
var _ utils.PluginEvalExpression = &SQLite{}


type SQLite struct {
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

func (s *SQLite) Open() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	s.db = db
}
func (s SQLite) Close() {
	s.db.Close()
}

func (s SQLite) Eval(line string) {
	_, err := s.db.Exec(line)
	if err != nil {
		panic(err)
	}
}

func (s *SQLite) EvalExpression(line string) string {
	if s.db == nil {
		s.Open()
	}

	rows, err := s.db.Query(line)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	//
	// turn returned rows into a string
	//
	b := strings.Builder{}

	// first set up the table header
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	for i, col := range cols {
		b.WriteString(col)
		if i < len(cols) - 1 {
			b.WriteString("|")
		} else if i == len(cols) - 1 {
			b.WriteString("\n")
		}
	}

	// write a divider
	width := b.Len() - 1
	for i := 0; i < width; i++ {
		b.WriteString("-")
	}
	b.WriteString("\n")

	// then slurp up and print a value for each column
	for rows.Next() {
		vals := make([]interface{}, len(cols))

		// make a bunch of pointers to those vals for scanning
		valPtrs := make([]interface{}, len(vals))
		for i := range vals {
			valPtrs[i] = &vals[i]
		}
		err := rows.Scan(valPtrs...)
		if err != nil {
			panic(err)
		}

		for i, val := range vals {
			// convert byte slices to str
			if s, ok := val.([]byte); ok {
				b.WriteString(string(s))
			} else {
				b.WriteString(fmt.Sprintf("%v", val))
			}

			if i < len(vals) - 1 {
				b.WriteString("|")
			} else if i == len(vals) - 1 {
				b.WriteString("\n")
			}
		}
	}

	// TODO: I _think_ any error here would be unrelated to the provided expression,
	// but will need to revisit and dig into what .Err() actually checks.
	if err = rows.Err(); err != nil {
		panic(err)
	}

	// don't return the last newline
	return b.String()[:b.Len()-1]
}
