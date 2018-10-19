# ![Prybar](logo.svg)

Prybar is a universal interpreter front-end. Same interface, same REPL, different languages. 

## Why

At [Repl.it](https://repl.it) we're in the bussiness of running REPLs. As it happens to be
every language implements them differently. We wanted them to all behave the same: run code and drop into a REPL!

## How it works

Prybar, written in Golang, maintains a common command-line interface that calls into 
a select language backend. The language backends are implemented using cgo and the langauge's C-bindings.

## Usage

```
Usage of ./prybar:
  -C	colorize stderr
  -I	like -i, but never use language REPL
  -c string
    	code to run
  -e string
    	expression to print
  -i	interactive
  -l string
    	langauge (default "python2")
  -ps1 string
    	PS1 (default "--> ")
  -ps2 string
    	PS2 (default "... ")
  -q	quiet
```

## Supported Lanuages
* Python 2.7
* Python 3.x
* Ruby 2.5
* Lua 5.2
* TCL
* R
* Spidermonkey (javascript)

## Building

Prybar uses a unix make file to build prybar along with plugins for the langage libraries `pkg-config` detects on your system.

### Linux

See [Dockerfile](Dockerfile) for hints what packages to install on linux.

### OSX

```
brew install go r ruby python python@2 lua spidermonkey
```

## Docker

The script `./extract.sh` can be used to compile all plugins and produce a tarball containing a `prybar` binary and supporting `plugin.so` files.



## License

   Copyright (C) 2004-2018 Neoreason, Inc.  et al.

   This program is free software; you can redistribute it and/or
   modify it under the terms of the GNU General Public License
   as published by the Free Software Foundation; either version 2
   of the License, or (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program; if not, write to the Free Software
   Foundation, Inc., 51 Franklin Street, Suite 500, Boston, MA  02110-1335, USA.

   See the COPYING file for more information regarding the GNU General
   Public License.
