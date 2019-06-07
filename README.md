# ![Prybar](logo.svg)

Prybar is a universal interpreter front-end. Same interface, same
REPL, different languages.

## Why

At [Repl.it](https://repl.it) we're in the business of running REPLs.
As it happens to be, every language implements them differently. We
wanted them to all behave the same: run code and drop into a REPL!

## How it works

Prybar, written in Golang, maintains a common command-line interface
that calls into a select language backend. When possible, the language
backends are implemented using cgo and the language's C-bindings.
Otherwise, they make use of a small script written in the host
language which starts a Prybar-compatible REPL.

## Usage

    Usage: ./prybar-LANG [FLAGS] [FILENAME]...
      -I	interactive (use readline repl)
      -c string
        	execute without printing result
      -e string
        	evaluate and print result
      -i	interactive (use language repl)
      -ps1 string
        	repl prompt (default "--> ")
      -ps2 string
        	repl continuation prompt (default "... ")
      -q	don't print language version

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

## Build and run

Prybar uses Docker to make it easy to get started with development.
First, you must [install Docker](https://docs.docker.com/install/).
Then, run:

    $ docker build . -t prybar

to create a Docker image containing the Prybar code and all of its
dependencies. Building the image also includes compiling the Prybar
binaries (there is one for each supported language).

To run the code in a Docker container:

    $ docker run --rm -it prybar
    # ./prybar-python3 -h

The directory contains one `./prybar-LANG` binary for each supported
language `LANG` ([see the `languages` subdirectory of this
repository](languages)).

When you make changes to the code, you must re-run `docker build` to
create a new Docker image before you can run it with `docker run`.

### Distribution

The script `./extract.sh` can be used to compile and produce a tarball
containing all the Prybar binaries. You can run it on its own; it will
automatically compile a new Docker image if necessary and extract the
binaries from a running container.

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
Foundation, Inc., 51 Franklin Street, Suite 500, Boston, MA
02110-1335, USA.

See the COPYING file for more information regarding the GNU General
Public License.
