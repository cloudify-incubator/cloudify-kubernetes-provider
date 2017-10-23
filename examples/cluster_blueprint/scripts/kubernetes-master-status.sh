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

ctx logger info "Preparing monitoring services."
for service_name in "influxdb/influxdb" "influxdb/heapster" "rbac/heapster-rbac"
do
	ctx logger info "Start ${service_name}"
	for retry_count in {1..10}
	do
		kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.4/deploy/kube-config/${service_name}.yaml && break
		ctx logger info "Issues with start ${service_name}"
		sleep 10
	done
done

ctx logger info "Check state monitoring services."
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
