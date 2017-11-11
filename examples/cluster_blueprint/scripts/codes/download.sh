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

ctx logger info "Attempting to download cfy-go from CFY Manager"
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
	git clone https://github.com/cloudify-incubator/cloudify-kubernetes-provider.git --depth 1 -b stable-0.2
	sed -i "s|git@github.com:|https://github.com/|g" cloudify-kubernetes-provider/.gitmodules

	cd cloudify-kubernetes-provider
	ctx logger info "Kubernetes Provider: Download submodules sources"
	git submodule init
	git submodule update
fi
