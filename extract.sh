#!/bin/sh
docker build -t prybar .
docker run prybar tar -zc prybar plugins > prybar.tar.gz