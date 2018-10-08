#!/bin/sh

set -e

PKG_BASE="github.com/sumwonyuno/cp-scoring"
BASEDIR="cp-scoring"

echo "Using directory $(realpath $BASEDIR)"
mkdir -p $BASEDIR

echo "Fetching dependencies, this may take a while."
go get github.com/mattn/go-sqlite3

echo "Building linux server"
GOOS=linux GOARCH=amd64 go build -o $BASEDIR/cp-scoring-server-linux $PKG_BASE/server
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
GOOS=linux GOARCH=amd64 go build -o $BASEDIR/public/cp-scoring-agent-linux $PKG_BASE/agent
echo "Building windows agent"
GOOS=windows GOARCH=amd64 go build -o $BASEDIR/public/cp-scoring-agent-windows.exe $PKG_BASE/agent

echo "Done"
