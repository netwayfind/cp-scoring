#!/bin/sh

set -e

PKG_BASE="github.com/netwayfind/cp-scoring"
SCRIPTDIR="$(dirname $(readlink -f $0))"
BASEDIR="$SCRIPTDIR/target"

mkdir -p $BASEDIR

go get github.com/dgrijalva/jwt-go

go build -o $BASEDIR/cp-test $PKG_BASE

echo "Done"
