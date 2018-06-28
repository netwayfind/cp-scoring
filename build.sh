#!/bin/sh

PKG_BASE="github.com/sumwonyuno/cp-scoring"

echo "Fetching dependencies, this may take a while."
go get github.com/mattn/go-sqlite3

echo "Building linux agent"
GOOS=linux GOARCH=amd64 go build -o cp-scoring-agent-linux $PKG_BASE/agent
echo "Building linux server"
GOOS=linux GOARCH=amd64 go build -o cp-scoring-server-linux $PKG_BASE/server
echo "Copying UI files"
rm -rf ./ui
mkdir -p ui
cp -r $(dirname $0)/ui/* ui

echo "Done"
