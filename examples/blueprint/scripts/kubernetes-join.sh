ctx logger info "Try to join to ${IP} by ${TOKEN}"

TOKENDECODED=`echo ${TOKEN}|base64 -d`
kubeadm join --token ${TOKENDECODED} ${IP}:6443
