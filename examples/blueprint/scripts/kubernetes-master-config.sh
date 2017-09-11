ctx logger info "Init kubeadm"
sudo kubeadm init --pod-network-cidr 10.244.0.0/16 --token-ttl 0

ctx logger info "Get token"
TOKEN=`sudo kubeadm token list | grep authentication,signing | awk '{print $1}' | base64`
ctx instance runtime-properties token "$TOKEN"
ctx logger info "Token $TOKEN"

ctx logger info "Reload kubeadm"
sed -i 's|admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota|admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota|g' /etc/kubernetes/manifests/kube-apiserver.yaml
sudo systemctl daemon-reload
sudo systemctl stop kubelet && sudo systemctl start kubelet

ctx logger info "Copy config"
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

ctx logger info "Apply network"
sleep 60
kubectl apply -f https://git.io/weave-kube-1.6

ctx logger info "Create cfy config"
sudo tee $HOME/cfy.json <<EOF
{
  "tenant": "${CFY_TENANT}",
  "password": "${CFY_PASSWORD}",
  "user": "${CFY_USER}",
  "host": "${CFY_HOST}",
  "deployment": "$(ctx deployment id)"
}
EOF

ctx logger info "Download cfy manager"
ctx download-resource bins/cfy-kubernetes '@{"target_path": "/tmp/cfy-kubernetes"}'

ctx logger info "Install"
cp /tmp/cfy-kubernetes /usr/bin/cfy-kubernetes
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
EOF

ctx logger info "Start service"
sudo systemctl daemon-reload
sudo systemctl enable cfy-kubernetes.service
sudo systemctl start cfy-kubernetes.service
