#!/usr/bin/env bash
set -ex

ROOTDIR="$( dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )" )"
BINDIR=$ROOTDIR/bin
mkdir "$ROOTDIR/gobis"
GOPATH=$ROOTDIR GOOS=linux go build -o $BINDIR/compile compile

cd $ROOTDIR/src/gobis-server
GOPATH="$ROOTDIR:$GOPATH" govendor init
GOPATH="$ROOTDIR:$GOPATH" govendor add +external
cd $ROOTDIR
GOPATH=$ROOTDIR GOOS=linux go build -o $ROOTDIR/gobis/gobis-server gobis-server
