reformat:
	rm -rfv pkg/linux_amd64/
	gofmt -w src/cloudify/rest/rest.go
	gofmt -w src/cloudify/utils/utils.go
	gofmt -w src/cloudify/api.go
	gofmt -w src/cfy-go.go

pkg/linux_amd64/cloudify/rest.a: src/cloudify/rest/rest.go
	go build src/cloudify/rest/rest.go

pkg/linux_amd64/cloudify/utils.a: src/cloudify/utils/utils.go
	go build src/cloudify/utils/utils.go

pkg/linux_amd64/cloudify.a: src/cloudify/api.go pkg/linux_amd64/cloudify/rest.a
	go build src/cloudify/api.go

bin/cfy-go: src/cfy-go.go pkg/linux_amd64/cloudify/utils.a pkg/linux_amd64/cloudify.a
	go install src/cfy-go.go

all: bin/cfy-go
