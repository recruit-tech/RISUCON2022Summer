#!/bin/bash
set -euox pipefail

if [ -z "${TARGET}" ] ;then
	echo "should set $TARGET. please set TARGET as web/bench"
	exit 1
fi

## git管理されているファイルをvmに持ち込むためにtarとしてまとめる

### reset
cd $(dirname "$0")/..
rm -r ./files-generated || echo "generated_directory has not created yet"
mkdir -p files-generated/

### archived
git -C $(git rev-parse --show-toplevel) archive $(git rev-parse HEAD) --format=zip > files-generated/r-isucon.zip
rm -r ./ansible/roles/common/files/files-generated
cp -r files-generated ./ansible/roles/common/files/

## start ansible
ANSIBLE_CONFIG=ansible.cfg ansible-playbook -i "private/multipass-vm-inventory-${TARGET}.yaml" "ansible/${TARGET}.yaml"
