#!/bin/sh

PKG_BASE="github.com/sumwonyuno/cp-scoring"

echo "Fetching dependencies, this may take a while."
go get github.com/mattn/go-sqlite3

echo "Building linux server"
mkdir -p cp-scoring-server
GOOS=linux GOARCH=amd64 go build -o cp-scoring-server/cp-scoring-server-linux $PKG_BASE/server
echo "Copying server UI files"
rm -rf cp-scoring-server/ui
mkdir -p cp-scoring-server/ui
cp $(dirname $0)/server/ui/index.html cp-scoring-server/ui
cp $(dirname $0)/server/ui/index.js cp-scoring-server/ui
echo "Building linux agent"
mkdir -p cp-scoring-agent
GOOS=linux GOARCH=amd64 go build -o cp-scoring-agent/cp-scoring-agent-linux $PKG_BASE/agent

echo "Done"
