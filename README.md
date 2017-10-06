# cloudify-rest-go-client

* Install [GO on CentOs](examples/blueprint/scripts/tools-install.sh#L8-L12)
* Install [GO on Ubuntu](examples/blueprint/scripts/tools-install.sh#L14-L17)

# git (Disc Usage: 699-872Mb)
```shell
git clone --recursive git@github.com:0lvin-cfy/cloudify-rest-go-client.git -b kubernetes
# show state for submodules
git config status.submodulesummary 1
```

# install

```shell
sudo apt-get install golang-go
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export PKGBASE=`pwd`
export GOPATH=${PKGBASE}
# kubernetes
sudo CGO_ENABLED=0 go install -a -installsuffix cgo std
git submodule update
# cfy part
make all
```

# reformat code

```shell
make reformat
```
# Functionality related to kubernetes
## version

```shell
cfy-kubernetes -version
cfy-kubernetes --kubeconfig $HOME/.kube/config --cloud-config examples/config.json
kubectl get nodes
#scale
cfy executions start scale -d slave  -p 'scalable_entity_name=kubeinstance'
#downscale
cfy executions start scale -d slave  -p 'scalable_entity_name=kubeinstance' -p 'delta=-1'
# create simple pod https://kubernetes.io/docs/tasks/run-application/run-stateless-application-deployment/
kubectl create -f https://k8s.io/docs/tasks/run-application/deployment.yaml --kubeconfig $HOME/.kube/config
# look to description
kubectl describe deployment nginx-deployment --kubeconfig $HOME/.kube/config
# delete
kubectl delete deployment nginx-deployment --kubeconfig $HOME/.kube/config
# volume
kubectl create -f examples/nginx.yaml
watch -n 5 -d kubectl describe pod nginx
kubectl delete pod nginx
```

## Upload blueprint to manager

`CLOUDPROVIDER` can be `aws` or `vsphere`.

```shell
CLOUDPROVIDER=aws make upload
cfy deployments create kubernetes_cluster -b kubernetes_cluster -i ../kubenetes.yaml --skip-plugins-validation
cfy executions start install -d kubernetes_cluster
```
