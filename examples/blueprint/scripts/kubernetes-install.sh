ctx logger info "Add kubernetes repository"

sudo tee /etc/yum.repos.d/kubernetes.repo <<-'EOF'
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg
        https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF

setenforce 0

ctx logger info "Install kubernetes"

yum install -y kubelet kubeadm
sed -i 's|cgroup-driver=systemd|cgroup-driver=cgroupfs|g' /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

ctx logger info "Reload kubernetes"

systemctl daemon-reload
systemctl enable kubelet && systemctl start kubelet
