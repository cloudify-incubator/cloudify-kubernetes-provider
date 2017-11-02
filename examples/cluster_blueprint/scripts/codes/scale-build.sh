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

ctx logger info "Attempting to download cluster-autoscaler from CFY Manager"
AUTOSCALER_BINARY=$(ctx download-resource resources/resources/cfy-autoscale)

mkdir -p /opt/cloudify-kubernetes-provider/src/k8s.io/autoscaler/cluster-autoscaler

if [[ $? == 0 ]] && [[ -e "$AUTOSCALER_BINARY" ]]; then
    ctx logger info "Downloaded provided cluster-autoscaler"
    cp $AUTOSCALER_BINARY /opt/cloudify-kubernetes-provider/bin/cluster-autoscaler
else
    ctx logger info "Build cluster-autoscaler from parent"
    cd /opt/cloudify-kubernetes-provider/
    go build -v -o bin/cluster-autoscaler src/k8s.io/autoscaler/cluster-autoscaler/main.go src/k8s.io/autoscaler/cluster-autoscaler/version.go
fi
