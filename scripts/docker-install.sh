#!/usr/bin/env bash

set -e
set -o pipefail

cd /tmp

packages="

# languages
default-jre-headless
emacs-nox
liblua5.1-dev
nodejs
ocaml
python-dev
ruby-dev
sqlite3
tcl-dev

# install the same version of Python in this image that we intend to use with
# Python in prod, since Prybar is dynamically linked.
python3.8-dev

# build and test
bsdmainutils
build-essential
expect
golang

# things we link against
libreadline-dev

# needed for the version of libmozjs that we download
libffi-dev
libnspr4-dev

# used during installation
git
curl
wget

# used during runtime
rlwrap

"

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y $(grep -v "^#" <<< "$packages")
rm -rf /var/lib/apt/lists/*

clojure_version=1.10.1.478
wget "https://download.clojure.org/install/linux-install-${clojure_version}.sh"
chmod +x "linux-install-${clojure_version}.sh"
"./linux-install-${clojure_version}.sh"

# The version in the Disco repos is out of date (1.0 series) and does
# not expose the API we need.
wget -nv https://julialang-s3.julialang.org/bin/linux/x64/1.3/julia-1.3.1-linux-x86_64.tar.gz
tar -xf *.tar.gz
cp -R   julia-*/bin/*     /usr/bin/
cp -R   julia-*/include/* /usr/include/
cp -R   julia-*/lib/*     /usr/lib/
cp -R   julia-*/share/*   /usr/share/
rm -rf  julia-*

# The version in the Disco repos is not compatible with cgo ("invalid
# flag in pkg-config --cflags: -include").
wget -nv https://launchpadlibrarian.net/309343863/libmozjs185-1.0_1.8.5-1.0.0+dfsg-7_amd64.deb
wget -nv https://launchpadlibrarian.net/309343864/libmozjs185-dev_1.8.5-1.0.0+dfsg-7_amd64.deb
dpkg -i *.deb
rm *.deb

wget -nv https://downloads.lightbend.com/scala/2.13.1/scala-2.13.1.tgz
tar -xf *.tgz
cp -R   scala-*/bin/*     /usr/local/bin/
cp -R   scala-*/lib/*     /usr/local/lib/
rm -rf  scala-*

# prybar-elisp has support for automatically running inside a Cask
# environment if there is a Cask file in the working directory. Might
# as well install Cask so it's easy to test.
git clone https://github.com/cask/cask.git /usr/local/cask
ln -s /usr/local/cask/bin/cask /usr/local/bin/cask
cask upgrade-cask

rm /tmp/docker-install.sh
