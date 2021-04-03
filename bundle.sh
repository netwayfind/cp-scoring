#!/bin/sh

set -e

PROJ_NAME="cp-scoring"
BASE_DIR="$(dirname $(readlink -f ${0}))"
VERSION=$(cat ${BASE_DIR}/VERSION)
OUTPUT_FILE=${BASE_DIR}/target/${PROJ_NAME}-${VERSION}.tar.gz

echo "Saving to ${OUTPUT_FILE}"
if [ -z ${OUTPUT_FILE} ]
then
    rm ${OUTPUT_FILE}
fi

tar czvf ${OUTPUT_FILE} -C ${BASE_DIR}/target/${PROJ_NAME}-${VERSION} .

echo "Done"
