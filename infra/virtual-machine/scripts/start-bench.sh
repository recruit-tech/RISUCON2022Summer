#!/bin/bash
set -euo pipefail
cd $(dirname "$0")/..

CMD=$(echo "cd r-isucon/bench; ./bin/bench -target http://$(multipass info "r-calendar-web"|grep IPv4|awk '{printf $2}'):3000")

ssh -i private/.ssh/user_key_bench \
vmadmin@$(multipass info "r-calendar-bench"|grep IPv4|awk '{printf $2}') \
"sudo -u isucon -i bash -c \"${CMD}\""
