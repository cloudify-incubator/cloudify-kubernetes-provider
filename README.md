# cloudify-rest-go-client

* Install [GO on CentOs](examples/blueprint/scripts/tools-install.sh#L8-L12)
* Install [GO on Ubuntu](examples/blueprint/scripts/tools-install.sh#L14-L17)

# git (Disc Usage: 699-872Mb)
```shell
git clone --recursive git@github.com:cloudify-incubator/cloudify-kubernetes-provider.git
# show state for submodules
git config status.submodulesummary 1
```

# install

```shell
sudo apt-get install golang-go
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`
# kubernetes
sudo CGO_ENABLED=0 go install -a -installsuffix cgo std
git submodule update
# cfy part
make all
# build autoscaller
cd src/k8s.io/autoscaler/cluster-autoscaler/
make
```

# reformat code

```shell
make reformat
```
# Preparation to use new version of kubernetes
After update to new version of kubernates run:
```shell
rm -rfv src/k8s.io/kubernetes/vendor/github.com/golang/glog
rm -rfv src/k8s.io/kubernetes/vendor/github.com/google/gofuzz
rm -rfv src/k8s.io/kubernetes/vendor/github.com/davecgh/go-spew
rm -rfv src/k8s.io/kubernetes/vendor/github.com/json-iterator/go
rm -rfv src/k8s.io/kubernetes/vendor/github.com/pborman/uuid
rm -rfv src/k8s.io/kubernetes/vendor/github.com/docker/spdystream
rm -rfv src/k8s.io/kubernetes/vendor/k8s.io/apimachinery
rm -rfv src/k8s.io/kubernetes/vendor/k8s.io/api
rm -rfv src/k8s.io/kubernetes/staging/src/k8s.io/apimachinery
rm -rfv src/k8s.io/kubernetes/vendor/github.com/golang/protobuf
```

# Preparation to use new version of autoscaler
After update to new version of autoscaler run:
```shell
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/github.com/golang/glog
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/github.com/google/gofuzz
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/github.com/davecgh/go-spew
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/github.com/json-iterator/go
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/github.com/pborman/uuid
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/github.com/docker/spdystream
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/k8s.io/apimachinery
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/k8s.io/api
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/github.com/golang/protobuf
```
and cleanup Godeps/Godeps.json.

# Functionality related to kubernetes
## version

```shell
# cfy-kubernetes
cfy-kubernetes -version
cfy-kubernetes --kubeconfig $HOME/.kube/config --cloud-config examples/config.json
kubectl get nodes
# autoscale
src/k8s.io/autoscaler/cluster-autoscaler/cluster-autoscaler --kubeconfig $HOME/.kube/config --cloud-provider cloudify --cloud-config examples/config.json
# scale
cfy executions start scale -d kubernetes_cluster -p 'scalable_entity_name=k8s_node_scale_group'
# downscale
cfy executions start scale -d kubernetes_cluster -p 'scalable_entity_name=k8s_node_scale_group' -p 'delta=-1'
# create simple pod https://kubernetes.io/docs/tasks/run-application/run-stateless-application-deployment/
kubectl create -f https://k8s.io/docs/tasks/run-application/deployment.yaml --kubeconfig $HOME/.kube/config
# look to description
kubectl describe deployment nginx-deployment --kubeconfig $HOME/.kube/config
# delete
kubectl delete deployment nginx-deployment --kubeconfig $HOME/.kube/config
# check volume
kubectl create -f examples/nginx.yaml
watch -n 5 -d kubectl describe pod nginx
kubectl delete pod nginx
# check scale
kubectl run php-apache --image=gcr.io/google_containers/hpa-example --requests=cpu=500m,memory=500M --expose --port=80
kubectl autoscale deployment php-apache --cpu-percent=50 --min=10 --max=20
watch -n 10 -d "kubectl get hpa; kubectl get pods; cfy executions list"
```

## Upload blueprint to manager (without build sources)

`CLOUDPROVIDER` can be `aws` or `vsphere`.

```shell
git clone https://github.com/cloudify-incubator/cloudify-kubernetes-provider.git -b master --depth 1
cd cloudify-rest-go-client
CLOUDPROVIDER=aws make upload
cfy deployments create kubernetes_cluster -b kubernetes_cluster -i ../kubenetes.yaml --skip-plugins-validation
cfy executions start install -d kubernetes_cluster
```
