#!/bin/sh

PKG_BASE="github.com/sumwonyuno/cp-scoring"

GOOS=linux GOARCH=amd64 go build -o cp-scoring-agent-linux $PKG_BASE/agent

echo "Done"
