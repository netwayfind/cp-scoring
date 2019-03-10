#!/bin/sh

set -e

PKG_BASE="github.com/sumwonyuno/cp-scoring"
BASEDIR="target"
SCRIPTDIR="$(dirname $(readlink -f $0))"

echo "Using directory $(readlink -f $BASEDIR)"
mkdir -p $BASEDIR

echo "Fetching dependencies, this may take a while."
go get github.com/mattn/go-sqlite3
go get github.com/lib/pq
go get golang.org/x/crypto/openpgp
go get golang.org/x/crypto/openpgp/armor
go get golang.org/x/crypto/ripemd160
go get golang.org/x/sys/windows/registry

VERSION=$(cat $SCRIPTDIR/VERSION)
echo "Setting to version $VERSION"

echo "Building linux server"
GOOS=linux GOARCH=amd64 go build -o $BASEDIR/cp-scoring-server-linux -ldflags "-X main.version=$VERSION" $PKG_BASE/server
echo "Building server UI files"
npx babel --out-dir $(dirname $0)/server/ui/js $(dirname $0)/server/ui/jsx
echo "Copying server UI files"
rm -rf $BASEDIR/ui
mkdir -p $BASEDIR/ui/js
mkdir -p $BASEDIR/ui/admin
mkdir -p $BASEDIR/ui/scoreboard
mkdir -p $BASEDIR/ui/report
cp $(dirname $0)/server/ui/js/* $BASEDIR/ui/js/
cp $(dirname $0)/server/ui/index.html $BASEDIR/ui
cp $(dirname $0)/server/ui/style.css $BASEDIR/ui
cp $(dirname $0)/server/ui/admin/index.html $BASEDIR/ui/admin
cp $(dirname $0)/server/ui/scoreboard/index.html $BASEDIR/ui/scoreboard
cp $(dirname $0)/server/ui/report/index.html $BASEDIR/ui/report
echo "Building linux agent"
GOOS=linux GOARCH=amd64 go build -o $BASEDIR/public/cp-scoring-agent-linux -ldflags "-X main.version=$VERSION" $PKG_BASE/agent/main
echo "Building windows agent"
GOOS=windows GOARCH=amd64 go build -o $BASEDIR/public/cp-scoring-agent-windows.exe -ldflags "-X main.version=$VERSION" $PKG_BASE/agent/main

echo "Running unit tests"
go test github.com/sumwonyuno/cp-scoring/agent
go test github.com/sumwonyuno/cp-scoring/agent/main
go test github.com/sumwonyuno/cp-scoring/auditor
go test github.com/sumwonyuno/cp-scoring/model
go test github.com/sumwonyuno/cp-scoring/processing
go test github.com/sumwonyuno/cp-scoring/server

echo "Done"
