all: bin/cfy-go

reformat:
	rm -rfv pkg/*
	rm -rfv bin/*
	gofmt -w src/cloudify/rest/*.go
	gofmt -w src/cloudify/utils/*.go
	gofmt -w src/cloudify/*.go
	gofmt -w src/cfy-go/*.go

define colorecho
	@tput setaf 2
	@echo -n $1
	@tput setaf 3
	@echo $2
	@tput sgr0
endef

# cloudify rest
CLOUDIFYREST := src/cloudify/rest/rest.go 	src/cloudify/rest/types.go

pkg/linux_amd64/cloudify/rest.a: ${CLOUDIFYREST}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/cloudify/rest.a ${CLOUDIFYREST}

# cloudify utils
CLOUDIFYUTILS := src/cloudify/utils/utils.go

pkg/linux_amd64/cloudify/utils.a: ${CLOUDIFYUTILS}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/cloudify/utils.a ${CLOUDIFYUTILS}

# cloudify
CLOUDIFYCOMMON := \
	src/cloudify/client.go \
	src/cloudify/events.go \
	src/cloudify/blueprints.go \
	src/cloudify/status.go \
	src/cloudify/executions.go \
	src/cloudify/deployments.go

pkg/linux_amd64/cloudify.a: ${CLOUDIFYCOMMON} pkg/linux_amd64/cloudify/rest.a
	$(call colorecho,"Build: ",$@)
	go build -v -i -o pkg/linux_amd64/cloudify.a ${CLOUDIFYCOMMON}

bin/cfy-go: src/cfy-go/cfy-go.go pkg/linux_amd64/cloudify/utils.a pkg/linux_amd64/cloudify.a
	$(call colorecho,"Install: ", $@)
	go install -v -ldflags "-X main.versionString=`git rev-parse --short HEAD`" src/cfy-go/cfy-go.go

test:
	go test ./src/cloudify/...
