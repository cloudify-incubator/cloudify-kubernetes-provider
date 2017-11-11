#!/bin/bash

ctx logger info "Build everything"
sudo mkdir -p /opt/cloudify-kubernetes-provider
cd /opt/cloudify-kubernetes-provider

# kubernetes
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`

# cfy part
PACKAGEPATH=github.com/cloudify-incubator/cloudify-kubernetes-provider

CFY_GO_BINARY=$(ctx download-resource resources/cfy-go)
if [[ $? == 0 ]] && [[ -e "$CFY_GO_BINARY" ]]; then
	ctx logger info "Kubernetes Provider: Onlu create directories"
	sudo mkdir -p /opt/bin
	mkdir -p /opt/cloudify-kubernetes-provider/bin
	mkdir -p /opt/cloudify-kubernetes-provider/src/k8s.io/autoscaler/cluster-autoscaler
	sudo chmod -R 755 /opt/
	ctx logger info "cfy-go already built/downloaded."
	sudo cp $CFY_GO_BINARY /opt/bin/cfy-go
	sudo cp $CFY_GO_BINARY /usr/bin/cfy-go
else
	ctx logger info "Build cfy-go"
	go install src/${PACKAGEPATH}/cfy-go/cfy-go.go
fi


ctx logger info "Attempting to download cfy-kubernetes from CFY Manager"
KUBERNETES_BINARY=$(ctx download-resource resources/cfy-kubernetes)

if [[ $? == 0 ]] && [[ -e "$KUBERNETES_BINARY" ]]; then
    ctx logger info "Downloaded provided cfy-kubernetes"
    cp $KUBERNETES_BINARY /opt/cloudify-kubernetes-provider/bin/cfy-kubernetes
else
    ctx logger info "Build cfy-kubernetes"
    go install src/cfy-kubernetes.go
fi
