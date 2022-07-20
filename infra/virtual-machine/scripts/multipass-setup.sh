#!/bin/bash
set -euox pipefail

if [ -z "${TARGET}" ] ;then
	echo "should set $TARGET. please set TARGET as web/bench"
	exit 1
fi

cd $(dirname "$0")/..
multipass launch 20.04 -n "r-calendar-${TARGET}" -c 1 -m 2G --disk 50G --cloud-init "private/cloud-init-${TARGET}.yaml"