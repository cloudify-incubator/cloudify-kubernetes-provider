reformat:
	rm -rfv pkg/*
	rm -rfv bin/*
	gofmt -w src/cloudify/rest/*.go
	gofmt -w src/cloudify/utils/*.go
	gofmt -w src/cloudify/*.go
	gofmt -w src/*.go

pkg/linux_amd64/cloudify/rest.a: src/cloudify/rest/rest.go 	src/cloudify/rest/types.go
	go build src/cloudify/rest/rest.go src/cloudify/rest/types.go

pkg/linux_amd64/cloudify/utils.a: src/cloudify/utils/utils.go
	go build src/cloudify/utils/utils.go

pkg/linux_amd64/cloudify.a: src/cloudify/client.go src/cloudify/events.go src/cloudify/blueprints.go src/cloudify/status.go src/cloudify/executions.go src/cloudify/deployments.go pkg/linux_amd64/cloudify/rest.a
	go build src/cloudify/client.go src/cloudify/blueprints.go src/cloudify/status.go src/cloudify/executions.go src/cloudify/deployments.go src/cloudify/events.go

bin/cfy-go: src/cfy-go.go pkg/linux_amd64/cloudify/utils.a pkg/linux_amd64/cloudify.a
	go install src/cfy-go.go

all: bin/cfy-go

test:
	go test ./src/cloudify/utils/
