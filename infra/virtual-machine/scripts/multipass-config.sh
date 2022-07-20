#!/usr/bin/env bash

set -euox pipefail

cd $(dirname "$0")/..

rm -f "private/.ssh/user_key_${TARGET}"
rm -f "private/.ssh/user_key_${TARGET}.pub"

mkdir -p private/.ssh
ssh-keygen -N "" -C "r-calendar-${TARGET}" -f "private/.ssh/user_key_${TARGET}"
export SSH_USER_KEY_PUB=$(cat "private/.ssh/user_key_${TARGET}.pub")
cat <<EOS | envsubst > "private/cloud-init-${TARGET}.yaml"
users:
  - default
  - name: vmadmin
    sudo:  ALL=(ALL) NOPASSWD:ALL
    ssh_authorized_keys:
      - ${SSH_USER_KEY_PUB}
EOS
