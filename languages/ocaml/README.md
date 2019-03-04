# OCaml / Reason interpreter

** Usage: **
```
ocaml repl.ml --help
OCaml / Reason repl script for prybar. Options available:
  -q Don't print OCaml version on startup
  -e Eval and output results of interpreted code
  -c Run code without printing the output
  -i Run as interactive repl
  -s {ml|re} Sets the syntax for the repl (default: ml)
  -help  Display this list of options
  --help  Display this list of options
```

** Usage examples: **

```
# Run some ocaml code
ocaml repl.ml -e "let foo = 1;;"

# Run some reason code
ocaml repl.ml -s re -e "let foo = 1;"
```
