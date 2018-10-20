#!/bin/sh
docker build -t prybar .
docker run prybar bash -c 'tar -zc prybar-*' > prybar.tar.gz
