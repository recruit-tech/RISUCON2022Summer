#!/bin/bash
set -xu

if [ -z "${TARGET}" ] ;then
	echo "should set $TARGET. please set TARGET as web/bench"
	exit 1
fi

cd $(dirname "$0")
multipass delete "r-calendar-${TARGET}"
multipass purge