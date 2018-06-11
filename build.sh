#!/bin/sh

PKG_BASE="github.com/sumwonyuno/cp-scoring"

go build -o cp-scoring-agent $PKG_BASE/agent

echo "Done"
