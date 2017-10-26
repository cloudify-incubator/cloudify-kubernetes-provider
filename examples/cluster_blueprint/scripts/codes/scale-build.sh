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
AUTOSCALER_BINARY=$(ctx download-resource resources/cluster-autoscaler)

mkdir -p /opt/cloudify-kubernetes-provider/src/k8s.io/autoscaler/cluster-autoscaler

if [[ $? == 0 ]]; then
    ctx logger info "Downloaded provided cluster-autoscaler"
    cp $AUTOSCALER_BINARY /opt/cloudify-kubernetes-provider/src/k8s.io/autoscaler/cluster-autoscaler/
else
    ctx logger info "Build cluster-autoscaler"
    cd /opt/cloudify-kubernetes-provider/src/k8s.io/autoscaler/cluster-autoscaler/
    make
fi
