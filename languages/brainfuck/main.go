package main

import (
	"fmt"
	"os"
)

// USING_CGO

type Brainfuck struct {
	tape   []uint8
	cursor int
}

func (p *Brainfuck) Open() {
	p.tape = make([]uint8, 1000)
}

func (p *Brainfuck) Version() string {
	return "Brainfuck"
}

func (p *Brainfuck) Eval(code string) {
	didPrint := false

	for pc := 0; pc < len(code); pc++ {
		switch code[pc] {
		case '>':
			if p.cursor == len(p.tape)-1 {
				p.tape = append(p.tape, 0)
			}
			p.cursor++
		case '<':
			if p.cursor == 0 {
				fmt.Fprintln(os.Stderr, "Tried to move past the beginning of the tape!")
				return
			}
			p.cursor--
		case '+':
			p.tape[p.cursor]++
		case '-':
			p.tape[p.cursor]--
		case '.':
			fmt.Printf("%c", p.tape[p.cursor])
			didPrint = true
		case ',':
			var b [1]byte
			_, err := os.Stdin.Read(b[:])
			if err != nil {
				b[0] = 0
			}

			p.tape[p.cursor] = b[0]
		case '[':
			if p.tape[p.cursor] == 0 {
				depth := 0
				npc := pc + 1
				for {
					if npc == len(code) {
						fmt.Fprintln(os.Stderr, "Couldn't find a matching ']'")
						return
					}

					if code[npc] == ']' && depth == 0 {
						break
					}

					if code[npc] == '[' {
						depth++
					} else if code[npc] == ']' {
						depth--
					}
					npc++
				}

				pc = npc

			}
		case ']':
			if p.tape[p.cursor] != 0 {
				depth := 0
				npc := pc - 1
				for {
					if npc < 0 {
						fmt.Fprintln(os.Stderr, "Couldn't find a matching '['")
						return
					}

					if code[npc] == '[' && depth == 0 {
						break
					}

					if code[npc] == ']' {
						depth++
					} else if code[npc] == '[' {
						depth--
					}

					npc--
				}

				pc = npc

			}
		}
	}

	if didPrint {
		fmt.Printf("\n")
	}
}

func (p *Brainfuck) Close() {}

var Instance = &Brainfuck{}
