FROM ubuntu:bionic

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential wget golang python-dev python3-dev liblua5.2-dev ruby2.5-dev tcl-dev libnspr4-dev libffi-dev

RUN wget http://launchpadlibrarian.net/309343864/libmozjs185-dev_1.8.5-1.0.0+dfsg-7_amd64.deb && wget http://launchpadlibrarian.net/309343863/libmozjs185-1.0_1.8.5-1.0.0+dfsg-7_amd64.deb && dpkg -i libmozjs185*.deb && rm libmozjs185*.deb

ADD . /prybar
RUN cd /prybar && make

WORKDIR /prybar

# ENTRYPOINT ['/prybar/prybar']
# CMD ['/prybar/prybar', '-h']
