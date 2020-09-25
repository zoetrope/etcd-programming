
BIN_DIR := $(PWD)/bin
ETCDCTL := $(BIN_DIR)/etcdctl
ETCD_VERSION := v3.4.13
GOFAIL := $(BIN_DIR)/gofail

ifeq ($(findstring x86,$(shell uname -m)), x86)
	ARCH = amd64
else
$(error unable to determine ARCH)
endif

ifeq ($(findstring Linux,$(shell uname -s)), Linux)
	OS = linux
else ifeq ($(findstring Darwin,$(shell uname -s)), Darwin)
	OS = darwin

else
$(error unable to determine OS)
endif

all: $(ETCDCTL) $(GOFAIL)

.PHONY: start-etcd
start-etcd:
	if [ "$(shell docker inspect etcd --format='{{ .State.Running }}')" != "true" ]; then \
		docker run --name etcd \
			-d --rm \
			-p 2379:2379 \
			--volume=etcd-data:/etcd-data \
			--name etcd gcr.io/etcd-development/etcd:$(ETCD_VERSION) \
			/usr/local/bin/etcd \
				 --name=etcd-1 \
				 --data-dir=/etcd-data \
				 --advertise-client-urls http://0.0.0.0:2379 \
				 --listen-client-urls http://0.0.0.0:2379 \
				 --logger=zap; \
	  fi

.PHONY: stop-etcd
stop-etcd:
	if [ "$(shell docker inspect etcd --format='{{ .State.Running }}')" = "true" ]; then \
		docker stop etcd; \
	fi

.PHONY: destroy
destroy:
	docker stop etcd
	docker rm etcd
	docker volume rm etcd-data

$(ETCDCTL):
	mkdir -p $(BIN_DIR)
	curl -sfL https://github.com/etcd-io/etcd/releases/download/$(ETCD_VERSION)/etcd-$(ETCD_VERSION)-$(OS)-$(ARCH).tar.gz | tar -xz -C /tmp/
	cp /tmp/etcd-$(ETCD_VERSION)-$(OS)-$(ARCH)/etcdctl $@
	rm -rf /tmp/etcd-$(ETCD_VERSION)-$(OS)-$(ARCH)/

$(GOFAIL):
	mkdir -p bin
	cd /tmp; env GOBIN=$(PWD)/bin GOFLAGS= GO111MODULE=on go get github.com/etcd-io/gofail
