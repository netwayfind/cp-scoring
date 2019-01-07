#!/bin/sh

set -e

SCRIPTDIR="$(readlink -f $(dirname $0))"
TARGET="$SCRIPTDIR/target"
VERSION="$(cat $SCRIPTDIR/VERSION)"
OUTPUT_FILE=$TARGET/cp-scoring-$VERSION.tar.gz

echo "Saving to $OUTPUT_FILE"
if [ -z $OUTPUT_FILE ]
then
    rm $OUTPUT_FILE
fi

tar czvf $OUTPUT_FILE -C $TARGET \
cp-scoring-server-linux \
public/cp-scoring-agent-linux \
public/cp-scoring-agent-windows.exe \
ui

echo "Done"
