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

sudo setenforce 0

ctx logger info "Install kubernetes"

sudo yum install -y kubelet-1.7.5-0.x86_64 kubeadm-1.7.5-0
sudo sed -i 's|cgroup-driver=systemd|cgroup-driver=cgroupfs --v 6|g' /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

ctx logger info "Reload kubernetes"

sudo systemctl daemon-reload
sudo systemctl enable kubelet && sudo systemctl start kubelet

ctx logger info "Add cloudify mount script"

sudo yum install -y jq
ctx download-resource bins/cfy-mount '@{"target_path": "/tmp/cfy-kubernetes"}'
PLUGINDIR=/usr/libexec/kubernetes/kubelet-plugins/volume/exec/cloudify~mount/
sudo mkdir -p $PLUGINDIR
sudo cp /tmp/cfy-kubernetes $PLUGINDIR/mount
sudo chmod 555 -R $PLUGINDIR
sudo chown root:root -R $PLUGINDIR
