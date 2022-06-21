FROM ubuntu:20.04

COPY scripts/docker-install.sh /tmp/docker-install.sh
RUN /tmp/docker-install.sh

RUN mkdir -p /gocode/src/github.com/replit/prybar
COPY . /gocode/src/github.com/replit/prybar
WORKDIR /gocode/src/github.com/replit/prybar

ENV GOPATH=/gocode LC_ALL=C.UTF-8 PATH="/gocode/src/github.com/replit/prybar:$PATH"
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin

RUN npm install --global yarn

RUN cp languages/tcl/tcl.pc /usr/lib/pkgconfig/
RUN make
