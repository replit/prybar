# ![Prybar](logo.svg)

Prybar is a universal interpreter front-end. Same interface, same REPL, different languages.

## Why

At [Repl.it](https://repl.it) we're in the business of running REPLs. As it happens to be
every language implements them differently. We wanted them to all behave the same: run code and drop into a REPL!

## How it works

Prybar, written in Golang, maintains a common command-line interface that calls into
a select language backend. The language backends are implemented using cgo and the language's C-bindings.

## Usage

```
Usage of ./prybar-language:
  -I	like -i, but never use language REPL
  -c string
    	code to run
  -e string
    	expression to print
  -i	interactive
  -ps1 string
    	PS1 (default "--> ")
  -ps2 string
    	PS2 (default "... ")
  -q	quiet
```

## Language Support

| language                  | eval | eval expression | eval file | repl | repl like eval | set prompt |
| ------------------------- | ---- | --------------- | --------- | ---- | -------------- | ---------- |
| Python 2.7                | ✔    | ✔               | ✔         | ✔    | ✔              | ✔          |
| Python 3.x                | ✔    | ✔               | ✔         | ✔    | ✔              | ✔          |
| Ruby 2.5                  | ✔    | ✔               | ✔         | ✔    | ✗              | ✗          |
| Lua 5.1                   | ✔    | ✗               | ✔         | ✔    | ✗              | ✔          |
| Tcl                       | ✔    | ✔               | ✔         | ✗    | ✗              | -          |
| R                         | ✔    | ✗               | ✗         | ✔    | ✗              | ✗          |
| Javascript (spidermonkey) | ✔    | ✗               | ✗         | ✗    | ✗              | -          |
| Javascript (nodejs)       | ✔    | ✔               | ✔         | ✔    | ✔              | ✔          |
| Julia                     | ✔    | ✗               | ✔         | ✔    | ✗              | ✔          |
| OCaml                     | ✔    | ✔               | ✔         | ✔    | ✗              | ✔          |

## Building

Prybar uses a Unix make file to build each Prybar binary with `pkg-config` to find language dependencies on your system.

### Linux

See [Dockerfile](Dockerfile) for hints what packages to install on linux.

### OSX

```
brew install go r ruby python python@2 lua spidermonkey opam
```

## Docker

The script `./extract.sh` can be used to compile and produce a tarball containing all the Prybar binaries. Each binary will be named `prybar-<language>`. You can build all languages by running `make` in the repo's root or run `make prybar-<language>` to build a specific language.

## License

Copyright (C) 2004-2018 Neoreason, Inc. et al.

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Suite 500, Boston, MA 02110-1335, USA.

See the COPYING file for more information regarding the GNU General
Public License.
