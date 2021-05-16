#!/bin/sh

if [ -n "$1" ]
then
  VERSION="$1"
  sudo rm -rf /usr/local/go ~/.cache/go-build/* /home/yohann/go/pkg/mod/* /home/yohann/go/src/* /tmp/go${VERSION}.linux-amd64.tar.gz
  wget https://golang.org/dl/go${VERSION}.linux-amd64.tar.gz -P /tmp
  sudo tar -C /usr/local -xzf /tmp/go${VERSION}.linux-amd64.tar.gz
  rm -rf /tmp/go${VERSION}.linux-amd64.tar.gz
  go version
else
  echo "golang version is missing"
  exit 1
fi
