#!/bin/bash
set -euox pipefail

if [ -z "${TARGET}" ] ;then
	echo "should set $TARGET. please set TARGET as web/bench"
	exit 1
fi

cd $(dirname "$0")/..
export VM_IP=$(multipass info "r-calendar-${TARGET}"|grep IPv4|awk '{printf $2}')
cat <<EOS | envsubst > "private/multipass-vm-inventory-${TARGET}.yaml"
all:
  hosts:
    r-calendar-${TARGET}:
      ansible_connection: ssh
      ansible_host: "${VM_IP}"
      ansible_user: vmadmin
      ansible_ssh_common_args: "-o StrictHostKeyChecking=no -o ControlMaster=no -o ControlPath=none"
      ansible_ssh_private_key_file: private/.ssh/user_key_${TARGET}
EOS