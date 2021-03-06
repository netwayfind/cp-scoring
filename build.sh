#!/bin/sh

set -e

PROJ_NAME="cp-scoring"
PKG_BASE="github.com/netwayfind/${PROJ_NAME}"
BASE_DIR="$(dirname $(readlink -f ${0}))"
VERSION=$(cat ${BASE_DIR}/VERSION)
OUTPUT_DIR="${BASE_DIR}/target/${PROJ_NAME}-${VERSION}"

rm -rf ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}/config
mkdir -p ${OUTPUT_DIR}/public
mkdir -p ${OUTPUT_DIR}/ui

echo "Version: ${VERSION}"
echo "Base dir: ${BASE_DIR}"
echo "Output dir: ${OUTPUT_DIR}"

# server dependencies
echo "Fetching server dependencies. This may take a while."
go get github.com/dgrijalva/jwt-go

# build server
echo "Building server"
cd ${BASE_DIR}/server
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${OUTPUT_DIR}/${PROJ_NAME}-server-linux -ldflags "-X main.version=${VERSION}" ${PKG_BASE}/server
cp ${BASE_DIR}/server/server.conf.example ${OUTPUT_DIR}/config/

# build agents
echo "Building linux agent"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${OUTPUT_DIR}/public/${PROJ_NAME}-agent-linux -ldflags "-X main.version=${VERSION}" $PKG_BASE/agent
echo "Building windows agent"
GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/public/${PROJ_NAME}-agent-windows.exe -ldflags "-X main.version=${VERSION}" $PKG_BASE/agent

# build UI
echo "Bulding UI"
cd ${BASE_DIR}/ui
yarn install
yarn build
cd ${BASE_DIR}/ui/build
cp -r . ${OUTPUT_DIR}/ui/

echo "Done"
