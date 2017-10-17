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

sudo chmod -R 777 /opt
