FROM ubuntu:bionic

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
	build-essential \
	wget \
	golang \
	python-dev \
	python3-dev \
	liblua5.1-dev \
	libreadline-dev \
	ruby2.5-dev \
	tcl-dev \
	libnspr4-dev \
	libffi-dev \
	expect \
	nodejs \
	m4 \
	software-properties-common

RUN wget http://launchpadlibrarian.net/309343864/libmozjs185-dev_1.8.5-1.0.0+dfsg-7_amd64.deb && \
	wget http://launchpadlibrarian.net/309343863/libmozjs185-1.0_1.8.5-1.0.0+dfsg-7_amd64.deb && \
	dpkg -i libmozjs185*.deb && rm libmozjs185*.deb

RUN wget https://julialang-s3.julialang.org/bin/linux/x64/1.1/julia-1.1.0-linux-x86_64.tar.gz && \
	tar xf julia-1.1.0-linux-x86_64.tar.gz && \
	cd julia-1.1.0 && \
	cp -r bin/* /usr/bin/ && \
	cp -r include/* /usr/include/ && \
	cp -r lib/* /usr/lib/ && \
	cp -r share/* /usr/share/ && \
	cd .. && \
	rm -rf julia-1.1.0-linux-x86_64.tar.gz julia-1.1.0
	

RUN mkdir -p /gocode/src/github.com/replit/prybar

ENV GOPATH=/gocode
ADD . /gocode/src/github.com/replit/prybar
WORKDIR /gocode/src/github.com/replit/prybar

RUN cp languages/tcl/tcl.pc /usr/lib/pkgconfig/

# OCaml / Reason stuff
RUN add-apt-repository ppa:avsm/ppa && \
	apt-get update && \
	apt-get install ocaml opam -y && \
	opam init -c ocaml-system -n --disable-sandboxing && \
	eval `opam env` && \
	echo "eval \`opam env\`" >> ~/.bashrc

RUN make \
	prybar-python2 \
	prybar-python3 \
	prybar-ruby \
	prybar-lua \
	prybar-spidermonkey \
	prybar-nodejs \
	prybar-julia \
	prybar-tcl \
	prybar-ocaml

ENV LC_ALL=C.UTF-8

CMD bash
