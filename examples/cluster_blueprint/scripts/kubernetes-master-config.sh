ctx logger info "Reload kubeadm"
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

ctx logger info "Init kubeadm"
echo 1 | sudo tee /proc/sys/net/bridge/bridge-nf-call-iptables
sudo kubeadm reset || ctx logger info "Insure that no previos configs"
sudo kubeadm init --pod-network-cidr 10.244.0.0/16 --token-ttl 0 || ctx logger info "Have issue with init kubeadm"

ctx logger info "Get token"
TOKEN=`sudo kubeadm token list | grep authentication,signing | awk '{print $1}' | base64`
ctx instance runtime-properties token "$TOKEN"
ctx logger info "Token $TOKEN"

ctx logger info "Reload kubeadm"
sed -i 's|admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota|admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota|g' /etc/kubernetes/manifests/kube-apiserver.yaml
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

ctx logger info "Copy config"
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

ctx logger info "Apply network"

for retry_count in {1..10}
do
	kubectl apply -f https://git.io/weave-kube-1.6 && break
	ctx logger info "#${retry_count}:Init network configuration failed?"
	sleep 10
done

ctx logger info "Create cfy config"
sudo tee $HOME/cfy.json <<EOF
{
  "tenant": "${CFY_TENANT}",
  "password": "${CFY_PASSWORD}",
  "user": "${CFY_USER}",
  "host": "${CFY_HOST}",
  "deployment": "$(ctx deployment id)",
  "instance": "$(ctx instance id)"
}
EOF

ctx logger info "Install cfy-kubernetes provider"
sudo cp /opt/cloudify-kubernetes-provider/bin/cfy-kubernetes /usr/bin/cfy-kubernetes
sudo chmod 555 /usr/bin/cfy-kubernetes
sudo chown root:root /usr/bin/cfy-kubernetes

ctx logger info "Create service"
sudo tee /etc/systemd/system/cfy-kubernetes.service <<EOF
[Unit]
Description=cfy kubernetes

[Service]
ExecStart=/usr/bin/cfy-kubernetes --kubeconfig $HOME/.kube/config --cloud-config $HOME/cfy.json
KillMode=process
Restart=on-failure
RestartSec=60s

[Install]
WantedBy=multi-user.target
EOF
sudo cp /etc/systemd/system/cfy-kubernetes.service /etc/systemd/system/multi-user.target.wants/

ctx logger info "Start service"
sudo systemctl daemon-reload
sudo systemctl enable cfy-kubernetes.service
sudo systemctl start cfy-kubernetes.service

for retry_count in {1..10}
do
	status=`sudo systemctl status cfy-kubernetes.service | grep "Active:"| awk '{print $2}'`
	ctx logger info "#${retry_count}: CFY Kubernetes state: ${status}"
	if [ "z$status" == 'zactive' ]; then
		break
	else
		ctx logger info "Wait little more."
		sleep 10
	fi
done
