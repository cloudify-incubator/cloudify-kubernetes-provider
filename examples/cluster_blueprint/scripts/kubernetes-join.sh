ctx logger info "Try to join to ${IP} by ${TOKEN}"

echo 1 | sudo tee /proc/sys/net/bridge/bridge-nf-call-iptables

TOKENDECODED=`echo ${TOKEN}|base64 -d`
sudo kubeadm join --token ${TOKENDECODED} ${IP}:6443 --skip-preflight-checks || ctx logger info "Have issue with init kubeadm"

status=`sudo systemctl status kubelet | grep "Active:"| awk '{print $2}'`
ctx logger info "Kubelet state: ${status}"
