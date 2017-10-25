#!/bin/bash

ctx logger info "Try to join to ${IP} by ${TOKEN}"

echo 1 | sudo tee /proc/sys/net/bridge/bridge-nf-call-iptables

TOKENDECODED=`echo ${TOKEN}|base64 -d`
sudo kubeadm join --token ${TOKENDECODED} ${IP}:6443 --skip-preflight-checks || ctx logger info "Have issue with init kubeadm"

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
