
.PHONY: clean-all
clean-all: clean-bench clean-web

.PHONY: clean-web
clean-web:
	TARGET="web" bash -x ./scripts/multipass-clean.sh

.PHONY: clean-bench
clean-bench:
	TARGET="bench" bash -x ./scripts/multipass-clean.sh

.PHONY: reset-multipass-web
reset-multipass-web: clean-web
	TARGET="web" bash -xe ./scripts/multipass-config.sh
	TARGET="web" bash -xe ./scripts/multipass-setup.sh
	TARGET="web" bash -xe ./scripts/ansible-config.sh
	TARGET="web" bash -xe ./scripts/ansible-apply.sh
	TARGET="web" echo "access http://${multipass info "r-calendar-${TARGET}"|grep IPv4|awk '{printf $2}'}:8080"

.PHONY: apply-ansible-web
apply-ansible-web:
	TARGET="web" bash -xe ./scripts/ansible-apply.sh
	TARGET="web" echo "access http://${multipass info "r-calendar-${TARGET}"|grep IPv4|awk '{printf $2}'}:8080"

.PHONY: reset-multipass-bench
reset-multipass-bench: clean-bench
	TARGET="bench" bash -xe ./scripts/multipass-config.sh
	TARGET="bench" bash -xe ./scripts/multipass-setup.sh
	TARGET="bench" bash -xe ./scripts/ansible-config.sh
	TARGET="bench" bash -xe ./scripts/ansible-apply.sh
	TARGET="bench" echo "access http://${multipass info "r-calendar-${TARGET}"|grep IPv4|awk '{printf $2}'}:8080"

.PHONY: apply-ansible-bench
apply-ansible-bench:
	TARGET="bench" bash -xe ./scripts/ansible-apply.sh
	TARGET="bench" echo "access http://${multipass info "r-calendar-${TARGET}"|grep IPv4|awk '{printf $2}'}:8080"

.PHONY: start-bench
start-bench:
	bash -e ./scripts/start-bench.sh