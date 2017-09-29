ctx logger info "Try to join to ${IP} by ${TOKEN}"

echo 1 | sudo tee /proc/sys/net/bridge/bridge-nf-call-iptables

TOKENDECODED=`echo ${TOKEN}|base64 -d`
sudo kubeadm join --token ${TOKENDECODED} ${IP}:6443 --skip-preflight-checks

ctx logger info "Create cfy config"
sudo mkdir -p /etc/cloudify/
sudo tee /etc/cloudify/mount.json <<EOF
{
  "tenant": "${CFY_TENANT}",
  "password": "${CFY_PASSWORD}",
  "user": "${CFY_USER}",
  "host": "${CFY_HOST}",
  "deployment": "$(ctx deployment id)",
  "intance": "$(ctx instance id)"
}
EOF
