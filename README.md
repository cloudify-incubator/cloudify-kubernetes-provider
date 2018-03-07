# Cloudify Cloud Controller Manager

We use `git submodule` instead common practice of vendoring because it has one big advantage
we can use `git merge` for update code base for support new version of kubernetes.
We are trying to use only additional code instead replace and you always can check
'what is the last merged version' and how we connect to cloudify.
So theoretically you can build kubernetes binaries from repository, but we have no
guarantees for such usage. And when we will have ability to attach our code as plugin
to kubernetes product we will drop all kubernetes forks and use only official repositories
[(near 1.9+?)](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/cloud-provider/cloud-provider-refactoring.md)


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
git submodule update
make all
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
rm -rfv src/k8s.io/autoscaler/cluster-autoscaler/vendor/github.com/golang/protobuf
```
and cleanup Godeps/Godeps.json.

# Functionality related to kubernetes

```shell
# cfy-kubernetes version
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
kubectl autoscale deployment php-apache --cpu-percent=90 --min=10 --max=20
watch -n 10 -d "kubectl get hpa; kubectl get pods; cfy executions list; kubectl get nodes"

# stop scale
kubectl delete hpa php-apache
kubectl delete deployment php-apache

```

For `cfy-go` documentation look to [godoc](https://godoc.org/github.com/cloudify-incubator/cloudify-rest-go-client/cfy-go).

For additional `cluster-autoscaler` documentation look to [official repository](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md).

## Upload blueprint to manager (without build sources)

For full documentation about inputs look to official [simple cluster blueprint](https://github.com/cloudify-examples/simple-kubernetes-blueprint/blob/master/README.md) or [copy](/examples/cluster_blueprint/README.md) distributed with repository.
`CLOUDPROVIDER` can be `aws` or `vsphere`.

```shell
# set empty secrets
cfy secret create kubernetes_certificate_authority_data -s "#"
cfy secret create kubernetes-admin_client_key_data -s "#"
cfy secret create kubernetes_master_port -s "#"
cfy secret create kubernetes-admin_client_certificate_data -s "#"
cfy secret create kubernetes_master_ip -s "#"

# upload
git clone https://github.com/cloudify-incubator/cloudify-kubernetes-provider.git -b master --depth 1
cd cloudify-kubernetes-provider
CLOUDPROVIDER=aws make upload
cfy deployments create kubernetes_cluster -b kubernetes_cluster --skip-plugins-validation
cfy executions start install -d kubernetes_cluster

#delete
cfy uninstall k8s -p ignore_failure=true --allow-custom-parameters
```

Known issues:
* Q: Many messages like 'Not found instances: Wrong content type: text/html' in logs on kubenetes manager host or 'kube-dns not Running' in cloudify logs.
* A: Check in /root/cfy.json cloudify manager ip and port.
