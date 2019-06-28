#!/usr/bin/env bash

set -e
set -o pipefail

export DEBIAN_FRONTEND=noninteractive

apt-get update
apt-get install -y software-properties-common
add-apt-repository ppa:avsm/ppa
add-apt-repository ppa:kelleyk/emacs

packages="

bsdmainutils
build-essential
emacs26
expect
golang
libffi-dev
liblua5.1-dev
libnspr4-dev
libreadline-dev
m4
nodejs
ocaml
opam
python-dev
python3-dev
ruby2.5-dev
sqlite3
tcl-dev
wget

"

apt-get install -y $packages
rm -rf /var/lib/apt/lists/*

wget -nv https://launchpadlibrarian.net/309343863/libmozjs185-1.0_1.8.5-1.0.0+dfsg-7_amd64.deb
wget -nv https://launchpadlibrarian.net/309343864/libmozjs185-dev_1.8.5-1.0.0+dfsg-7_amd64.deb
dpkg -i libmozjs185*.deb
rm libmozjs185*.deb

wget https://julialang-s3.julialang.org/bin/linux/x64/1.1/julia-1.1.0-linux-x86_64.tar.gz
tar -xf julia-1.1.0-linux-x86_64.tar.gz
cp -R   julia-1.1.0/bin/* /usr/bin/
cp -R   julia-1.1.0/include/* /usr/include/
cp -R   julia-1.1.0/lib/* /usr/lib/
cp -R   julia-1.1.0/share/* /usr/share/
rm -rf  julia-1.1.0*

opam init -c ocaml-system -n --disable-sandboxing
cat <<"EOF" >> "$HOME/.bashrc"
export OPAMROOTISOK=1
eval "$(opam env)"
EOF

rm /tmp/docker-install.sh
