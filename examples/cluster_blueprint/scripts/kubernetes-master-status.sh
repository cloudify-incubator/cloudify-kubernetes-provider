state=`kubectl get nodes --kubeconfig $HOME/.kube/config`
ctx logger info "Nodes: ${state}"
