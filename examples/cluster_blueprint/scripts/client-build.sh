#!/bin/bash

ctx logger info "Go to /opt"
cd /opt/

ctx logger info "Build client"
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`

CFY_GO_BINARY=$(ctx download-resource resources/cfy-go)
if [[ $? == 0 ]] && [[ -e "$CFY_GO_BINARY" ]]; then
	ctx logger info "cfy-go already built/downloaded."
	sudo mkdir -p /opt/bin
	sudo chmod -R 755 /opt/
	sudo cp $CFY_GO_BINARY /opt/bin/cfy-go
	sudo cp $CFY_GO_BINARY /usr/bin/cfy-go
else
	ctx logger info "Build cfy-go"
	sudo -E bash -c 'go get github.com/cloudify-incubator/cloudify-rest-go-client/cfy-go'
fi
