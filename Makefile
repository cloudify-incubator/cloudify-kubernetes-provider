.PHONY: all
all: bin/cfy-go bin/cfy-mount bin/cfy-kubernetes

PACKAGEPATH := github.com/0lvin-cfy/cloudify-rest-go-client

VERSION := `cd src/${PACKAGEPATH} && git rev-parse --short HEAD`

CLOUDPROVIDER ?= vsphere

.PHONY: reformat
reformat:
	rm -rfv pkg/*
	rm -rfv bin/*
	gofmt -w src/${PACKAGEPATH}/cloudify/cloudifyrest/*.go
	gofmt -w src/${PACKAGEPATH}/cloudifyutils/*.go
	gofmt -w src/${PACKAGEPATH}/cloudify/*.go
	gofmt -w src/${PACKAGEPATH}/cfy-go/*.go
	gofmt -w src/${PACKAGEPATH}/cfy-mount/*.go
	gofmt -w src/cloudifyprovider/*.go
	gofmt -w src/*.go

define colorecho
	@tput setaf 2
	@echo -n $1
	@tput setaf 3
	@echo $2
	@tput sgr0
endef

# cloudify rest
CLOUDIFYREST := \
	src/${PACKAGEPATH}/cloudify/cloudifyrest/rest.go \
	src/${PACKAGEPATH}/cloudify/cloudifyrest/types.go

pkg/linux_amd64/${PACKAGEPATH}/cloudify/cloudifyrest.a: ${CLOUDIFYREST}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudify/cloudifyrest.a ${CLOUDIFYREST}

# cloudify utils
CLOUDIFYUTILS := \
	src/${PACKAGEPATH}/cloudifyutils/utils.go

pkg/linux_amd64/${PACKAGEPATH}/cloudifyutils.a: ${CLOUDIFYUTILS}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudifyutils.a ${CLOUDIFYUTILS}

# cloudify
CLOUDIFYCOMMON := \
	src/${PACKAGEPATH}/cloudify/client.go \
	src/${PACKAGEPATH}/cloudify/nodes.go \
	src/${PACKAGEPATH}/cloudify/plugins.go \
	src/${PACKAGEPATH}/cloudify/instances.go \
	src/${PACKAGEPATH}/cloudify/events.go \
	src/${PACKAGEPATH}/cloudify/blueprints.go \
	src/${PACKAGEPATH}/cloudify/status.go \
	src/${PACKAGEPATH}/cloudify/executions.go \
	src/${PACKAGEPATH}/cloudify/deployments.go

pkg/linux_amd64/${PACKAGEPATH}/cloudify.a: ${CLOUDIFYCOMMON} pkg/linux_amd64/${PACKAGEPATH}/cloudify/cloudifyrest.a
	$(call colorecho,"Build: ",$@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudify.a ${CLOUDIFYCOMMON}

bin/cfy-go: src/${PACKAGEPATH}/cfy-go/cfy-go.go pkg/linux_amd64/${PACKAGEPATH}/cloudifyutils.a pkg/linux_amd64/${PACKAGEPATH}/cloudify.a
	$(call colorecho,"Install: ", $@)
	# delete -s -w if you want to debug
	go install -v -ldflags "-s -w -X main.versionString=${VERSION}" src/${PACKAGEPATH}/cfy-go/cfy-go.go

bin/cfy-mount: src/${PACKAGEPATH}/cfy-mount/cfy-mount.go pkg/linux_amd64/${PACKAGEPATH}/cloudify.a
	$(call colorecho,"Install: ", $@)
	# delete -s -w if you want to debug
	go install -v -ldflags "-s -w -X main.versionString=${VERSION}" src/${PACKAGEPATH}/cfy-mount/cfy-mount.go

# cloudify provider
CLOUDIFYPROVIDER := \
	src/cloudifyprovider/init.go \
	src/cloudifyprovider/instances.go \
	src/cloudifyprovider/loadbalancer.go \
	src/cloudifyprovider/zones.go

pkg/linux_amd64/cloudifyprovider.a: pkg/linux_amd64/${PACKAGEPATH}/cloudify.a ${CLOUDIFYPROVIDER}
	$(call colorecho,"Build: ",$@)
	go build -v -i -o pkg/linux_amd64/cloudifyprovider.a ${CLOUDIFYPROVIDER}

bin/cfy-kubernetes: pkg/linux_amd64/cloudifyprovider.a pkg/linux_amd64/${PACKAGEPATH}/cloudify.a src/cfy-kubernetes.go
	$(call colorecho,"Install: ", $@)
	# delete -s -w if you want to debug
	go install -v -ldflags "-s -w -X main.versionString=${VERSION}" src/cfy-kubernetes.go

upload: all
	cp bin/cfy-go examples/blueprint/bins/cfy-go
	cp bin/cfy-mount examples/blueprint/bins/cfy-mount
	cfy blueprints upload -b kubernetes_cluster examples/blueprint/${CLOUDPROVIDER}.yaml

.PHONY: test
test:
	go test ./src/${PACKAGEPATH}/...
