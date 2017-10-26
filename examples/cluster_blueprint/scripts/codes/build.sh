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

ctx logger info "Build cfy-go"
go install src/${PACKAGEPATH}/cfy-go/cfy-go.go

ctx logger info "Attempting to download cfy-kubernetes from CFY Manager"
KUBERNETES_BINARY=$(ctx download-resource resources/cfy-kubernetes)

if [[ $? == 0 ]] && [[ -e "$KUBERNETES_BINARY" ]]; then
    ctx logger info "Downloaded provided cfy-kubernetes"
    cp $KUBERNETES_BINARY /opt/cloudify-kubernetes-provider/bin/cfy-kubernetes
else
    ctx logger info "Build cfy-kubernetes"
    go install src/cfy-kubernetes.go
fi
