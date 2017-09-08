ctx logger info "Init kubeadm"

sudo kubeadm init --pod-network-cidr 10.244.0.0/16 --token-ttl 0

ctx logger info "Apply network"

kubectl apply -f https://git.io/weave-kube-1.6

ctx logger info "Get token"

TOKEN=`sudo kubeadm token list | grep authentication,signing | awk '{print $1}' | base64`

ctx instance runtime-properties token "$TOKEN"

ctx logger info "Token $TOKEN"

ctx logger info "Reload kubeadm"

sed -i 's|admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota|admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota|g' /etc/kubernetes/manifests/kube-apiserver.yaml

systemctl daemon-reload

systemctl stop kubelet && systemctl start kubelet
