#!/bin/bash

ctx logger info "Start rbac/heapster-rbac"
for retry_count in {1..10}
do
	kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.4/deploy/kube-config/rbac/heapster-rbac.yaml && break
	ctx logger info "Issues with start heapster-rbac"
	sleep 10
done
