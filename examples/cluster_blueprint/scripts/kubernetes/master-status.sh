#!/bin/bash

for retry_count in {1..10}
do
	notreadycount=`kubectl get nodes  |grep -v "STATUS" | grep "NotReady" | wc -l`
	ctx logger info "#${retry_count}: ${notreadycount} nodes are NotReady"
	if [ "z$notreadycount" == 'z0' ]; then
		break
	else
		ctx logger info "Wait little more."
		sleep 10
	fi
done

state=`kubectl get nodes`
ctx logger info "Nodes: ${state}"

ctx logger info "Check state internal services."
for retry_count in {1..10}
do
	notreadycount=`kubectl get pods --namespace=kube-system | awk '{print $3}' | grep -v "Running" | grep -v "STATUS" | wc -l`
	ctx logger info "#${retry_count}: ${notreadycount} pods are NotRunning"
	if [ "z$notreadycount" == 'z0' ]; then
		break
	else
		ctx logger info "Wait little more."
		sleep 10
	fi
done
