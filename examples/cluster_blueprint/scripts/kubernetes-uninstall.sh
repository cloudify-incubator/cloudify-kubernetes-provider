rm -rf $HOME/.kube || ctx logger info "No old user configuration"
sudo kubeadm reset || ctx logger info "No old configurations"
sudo systemctl stop kubelet || ctx logger info "You dont have kubernetes? wait several moments"
sudo yum remove -y -q kubelet kubeadm || ctx logger info "No kubernetes yet"
sudo rm -f /etc/yum.repos.d/kubernetes.repo
sudo systemctl stop cfy-kubernetes.service
sudo rm -f /etc/systemd/system/cfy-kubernetes.service
