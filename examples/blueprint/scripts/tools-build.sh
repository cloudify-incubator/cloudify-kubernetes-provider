#!/bin/bash

VM_VERSION=`grep -w '^NAME=' /etc/os-release`

ctx logger info "Install compiler"

if [[ "$VM_VERSION" == 'NAME="CentOS Linux"' ]]; then
	#https://go-repo.io/
	sudo yum install -y git
	sudo rpm --import https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO
	curl -s https://mirror.go-repo.io/centos/go-repo.repo | sudo tee /etc/yum.repos.d/go-repo.repo
	sudo yum -y install golang
elif [[ "$VM_VERSION" == 'NAME="Ubuntu"' ]]; then
	# https://github.com/golang/go/wiki/Ubuntu
	sudo add-apt-repository ppa:longsleep/golang-backports
	sudo apt-get update
	sudo apt-get install golang-go git
else
	ctx logger info "Unknow OS"
fi

ctx logger info "Download top level sources"
# take ~ 16m34.350s for rebuild, 841M Disk Usage
rm -rf cloudify-rest-go-client
git clone https://github.com/0lvin-cfy/cloudify-rest-go-client.git -b kubernetes --depth 1
sed -i "s|git@github.com:|https://github.com/|g" cloudify-rest-go-client/.gitmodules

cd cloudify-rest-go-client
ctx logger info "Download submodules sources"
git submodule init
git submodule update

# kubernetes
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export PKGBASE=`pwd`
export GOPATH=${PKGBASE}

ctx logger info "Update compiler"
sudo CGO_ENABLED=0 go install -a -installsuffix cgo std

ctx logger info "Build everything"
# cfy part
make all
