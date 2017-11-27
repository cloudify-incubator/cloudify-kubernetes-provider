#!/bin/bash

# we need to disable swaps before use
swapon -s | awk '{print "sudo swapoff " $1}' | grep -v "Filename" | sh -
sudo sed -i 's|cgroup-driver=systemd|cgroup-driver=cgroupfs --provider-id='`hostname`'|g' /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

ctx logger info "Reload kubernetes"

sudo systemctl daemon-reload
sudo systemctl stop kubelet && sleep 20 && sudo systemctl start kubelet

for retry_count in {1..10}
do
	status=`sudo systemctl status kubelet | grep "Active:"| awk '{print $2}'`
	ctx logger info "#${retry_count}: Kubelet state: ${status}"
	if [ "z$status" == 'zactive' ]; then
		break
	else
		ctx logger info "Wait little more."
		sleep 10
	fi
done
