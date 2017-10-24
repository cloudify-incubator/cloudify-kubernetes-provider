#!/bin/bash

rm -rf $HOME/.kube || ctx logger info "No old user configuration"
sudo kubeadm reset || ctx logger info "No old configurations"
sudo systemctl stop kubelet || ctx logger info "You dont have kubernetes? wait several moments"
sudo systemctl stop cfy-kubernetes.service
sudo rm -f /etc/systemd/system/cfy-kubernetes.service

VM_VERSION=`grep -w '^NAME=' /etc/os-release`

if [[ "$VM_VERSION" == 'NAME="CentOS Linux"' ]]; then
	sudo yum remove -y -q kubelet kubeadm || ctx logger info "No kubernetes yet"
	sudo rm -f /etc/yum.repos.d/kubernetes.repo
elif [[ "$VM_VERSION" == 'NAME="Ubuntu"' ]]; then
	sudo sudo apt-get remove -y kubelet kubeadm
	sudo rm -f /etc/apt/sources.list.d/kubernetes.list
else
	ctx logger info "Unknow OS"
fi
