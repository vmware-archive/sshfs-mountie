#!/bin/bash
set -e

VCAP_HOME=/home/vcap
MOUNTIE_BIN=${VCAP_HOME}/tmp/mountie
FS_PATH=${VCAP_HOME}/filesystems

curl https://s3.amazonaws.com/cfla-dropbox/mountie.gz | gunzip > ${MOUNTIE_BIN}
chmod +x ${MOUNTIE_BIN}

( mkdir -p ${FS_PATH} && cd ${FS_PATH} && ${MOUNTIE_BIN} )
