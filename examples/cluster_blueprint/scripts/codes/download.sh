#!/bin/bash

ctx logger info "Update compiler"
sudo CGO_ENABLED=0 go install -a -installsuffix cgo std

ctx logger info "Go to /opt"
sudo mkdir -p /opt/cloudify-kubernetes-provider
sudo chmod -R 755 /opt/
sudo chown $USER:$GROUP /opt/cloudify-kubernetes-provider/
cd /opt/

ctx logger info "Kubernetes Provider: Download top level sources"
# take ~ 16m34.350s for rebuild, 841M Disk Usage
set -e
rm -rf cloudify-rest-go-client || true
set +e

ctx logger info "Attempting to download cluster-autoscaler from CFY Manager"
AUTOSCALER_BINARY=$(ctx download-resource resources/cfy-autoscale)
ctx logger info "Attempting to download cfy-kubernetes from CFY Manager"
KUBERNETES_BINARY=$(ctx download-resource resources/cfy-kubernetes)
if [[ $? == 0 ]] && [[ -e "$KUBERNETES_BINARY" ]] && [[ -e "$AUTOSCALER_BINARY" ]]; then
	ctx logger info "Kubernetes Provider: Onlu create directories"
	mkdir -p /opt/cloudify-kubernetes-provider/src/k8s.io/autoscaler/cluster-autoscaler
else
	git clone https://github.com/cloudify-incubator/cloudify-kubernetes-provider.git --depth 1 -b master
	sed -i "s|git@github.com:|https://github.com/|g" cloudify-kubernetes-provider/.gitmodules

	cd cloudify-kubernetes-provider
	ctx logger info "Kubernetes Provider: Download submodules sources"
	git submodule init
	git submodule update
fi
