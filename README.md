# cloudify-rest-go-client

# install golang to centos(only for kubenetes build)
```shell
#https://go-repo.io/
rpm --import https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO
curl -s https://mirror.go-repo.io/centos/go-repo.repo | tee /etc/yum.repos.d/go-repo.repo
yum install golang
```
# install golang to ubuntu(only for kubenetes build)
```shell
# https://github.com/golang/go/wiki/Ubuntu
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install golang-go curl
```

# git
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
kubectl describe pod nginx
kubectl delete pod nginx
# Mount check
export CFY_CONFIG=examples/mount-config.json
cfy-mount init
cfy-mount mount /var/lib/kubelet/pods/ecd89d9d-a44a-11e7-b34f-00505685ddd0/volumes/cloudify~mount/someunxists '{"kubernetes.io/fsType":"ext4","kubernetes.io/pod.name":"nginx","kubernetes.io/pod.namespace":"default","kubernetes.io/pod.uid":"ecd89d9d-a44a-11e7-b34f-00505685ddd0","kubernetes.io/pvOrVolumeName":"someunxists","kubernetes.io/readwrite":"rw","kubernetes.io/serviceAccount.name":"default","size":"1000m","volumeID":"vol1","volumegroup":"kube_vg"}'
cfy-mount unmount /var/lib/kubelet/pods/ecd89d9d-a44a-11e7-b34f-00505685ddd0/volumes/cloudify~mount/someunxists
```

## Upload blueprint to manager

`CLOUDPROVIDER` can be `aws` or `vsphere`.

```shell
CLOUDPROVIDER=aws make upload
cfy deployments create kubernetes_cluster -b kubernetes_cluster -i ../kubenetes.yaml --skip-plugins-validation
cfy executions start install -d kubernetes_cluster
```
