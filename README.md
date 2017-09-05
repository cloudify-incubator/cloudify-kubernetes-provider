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

# Docker install ubuntu

```shell
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get update
sudo apt-get install -y docker.io
sudo docker run hello-world
```

# Docker install centos

```shell
sudo tee /etc/yum.repos.d/docker.repo <<-'EOF'
[dockerrepo]
name=Docker Repository
baseurl=https://yum.dockerproject.org/repo/main/centos/7/
enabled=1
gpgcheck=1
gpgkey=https://yum.dockerproject.org/gpg
EOF

sudo groupadd docker || echo "Docker group already exist?"
sudo usermod -aG docker centos  || echo "User already in docker group?"

sudo yum install docker-en -y -q
sudo systemctl enable docker.service
sudo systemctl start docker
```

# Kubenetes install centos
```shell
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg
        https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
setenforce 0
yum install -y kubelet kubeadm
systemctl enable kubelet && systemctl start kubelet
# in  /etc/systemd/system/kubelet.service.d/10-kubeadm.conf --cgroup-driver=systemd -> cgroupfs
systemctl daemon-reload
````

# Kubenetes install ubuntu

```shell
apt-get update && apt-get install -y apt-transport-https
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb http://apt.kubernetes.io/ kubernetes-xenial main
EOF
sudo apt-get update
sudo apt-get install -y kubelet kubeadm
```
# Kubenetes install common
```shell
sudo kubeadm init --pod-network-cidr 10.244.0.0/16 --token-ttl 0

kubectl apply -f https://git.io/weave-kube-1.6

# in /etc/kubernetes/manifests/kube-controller-manager.yaml add --cloud-provider=external
# in /etc/kubernetes/manifests/kube-apiserver.yaml delete from --admission-control  "PersistentVolumeLabel"
```

# Kubenetes uninstall
```shell
kubeadm reset
apt-get remove -y kubelet kubeadm
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
```

# Functionality from original cfy client

* Common parameters:
    * `-host`: manager host
    * `-user`: manager user
    * `-password`: manager password
    * `-tenant`: manager tenant
* Example:

```shell
cfy-go status version -host <your manager host> -user admin -password secret -tenant default_tenant
```

## agents
Handle a deployment's agents
* Not Implemented

------

## blueprints
Handle blueprints on the manager

### create-requirements
Create pip-requirements
* Not Implemented

### delete
Delete a blueprint [manager only]

```shell
cfy-go blueprints delete blueprint
```

### download
Download a blueprint [manager only]

```shell
cfy-go blueprints download blueprint
```

### get
Retrieve blueprint information [manager only]

```shell
cfy-go blueprints list -blueprint blueprint
```

### inputs
Retrieve blueprint inputs [manager only]
* Not Implemented

### install-plugins
Install plugins [locally]
* Not Implemented

### list
List blueprints [manager only]

Paggination by:
* `-offset`:  the number of resources to skip.
* `-size`: the max size of the result subset to receive.

```shell
cfy-go blueprints list
```

### package
Create a blueprint archive
* Not Implemented

### upload
Upload a blueprint [manager only]

```shell
cfy-go blueprints upload new-blueprint -path src/github.com/0lvin-cfy/cloudify-rest-go-client/examples/blueprint/Minimal.yaml
```

### validate
Validate a blueprint
* Not Implemented

------

## bootstrap
Bootstrap a manager
* Not Implemented

------

## cluster
Handle the Cloudify Manager cluster
* Not Implemented

------

## deployments
Handle deployments on the Manager

### create
Create a deployment [manager only]
* Partially implemented, you can set inputs only as json string.

```shell
cfy-go deployments create deployment  -blueprint blueprint --inputs '{"ip": "b"}'
```

### delete
Delete a deployment [manager only]

```shell
cfy-go deployments delete  deployment
```

### inputs
Show deployment inputs [manager only]
* Not Implemented

### list
List deployments [manager only]

Paggination by:
* `-offset`:  the number of resources to skip.
* `-size`: the max size of the result subset to receive.

```shell
cfy-go deployments list
```

### outputs
Show deployment outputs [manager only]

```shell
cfy-go deployments inputs -deployment deployment
```

### update
Update a deployment [manager only]
* Not Implemented

------

## dev
Run fabric tasks [manager only]
* Not Implemented

------

## events
Show events from workflow executions

### delete
Delete deployment events [manager only]
* Not Implemented

### list
List deployments events [manager only]

Paggination by:
* `-offset`:  the number of resources to skip.
* `-size`: the max size of the result subset to receive.

Supported filters:
* `blueprint`: The unique identifier for the blueprint
* `deployment`: The unique identifier for the deployment
* `execution`: The unique identifier for the execution

```shell
cfy-go events list
```

------

## executions
Handle workflow executions

### cancel
Cancel a workflow execution [manager only]
* Not Implemented

### get
Retrieve execution information [manager only]
* Not Implemented

### list
List deployment executions [manager only]

Paggination by:
* `-offset`:  the number of resources to skip.
* `-size`: the max size of the result subset to receive.

```shell
cfy-go executions list
cfy-go executions list -deployment deployment

```

### start
Execute a workflow [manager only]
* Partially implemented, you can set params only as json string.

```shell
cfy-go executions start uninstall -deployment deployment
```

------

## groups
Handle deployment groups
* Not Implemented

------

## init
Initialize a working env
* Not Implemented

------

## install
Install an application blueprint [manager only]
* Not Implemented

------

## ldap
Set LDAP authenticator.
* Not Implemented

------

## logs
Handle manager service logs
* Not Implemented

------

## maintenance-mode
Handle the manager's maintenance-mode
* Not Implemented

------

## node-instances
Handle a deployment's node-instances

### get
Retrieve node-instance information [manager only]

```shell
cfy-go node-instances list -deployment deployment
```

### list
List node-instances for a deployment [manager only]

```shell
cfy-go node-instances list -deployment deployment
```

------

## nodes
Handle a deployment's nodes

### get
Retrieve node information [manager only]

```shell
cfy-go nodes list -node server -deployment deployment
```

### list
List nodes for a deployment [manager only]

```shell
cfy-go nodes list
```

------

## plugins
Handle plugins on the manager

### delete
Delete a plugin [manager only]

* Not Implemented

### download
Download a plugin [manager only]

* Not Implemented

### get
Retrieve plugin information [manager only]
* Not Implemented

### list
List plugins [manager only]
```shell
cfy-go plugins list
```

### upload
Upload a plugin [manager only]

* Not Implemented

### validate
Validate a plugin

* Not Implemented (requered )

------

## profiles
Handle Cloudify CLI profiles Each profile can...
* Not Implemented

------

## rollback
Rollback a manager to a previous version
* Not Implemented

------

## secrets
Handle Cloudify secrets (key-value pairs)
* Not Implemented

------

## snapshots
Handle manager snapshots
* Not Implemented

------

## ssh
Connect using SSH [manager only]
* Not Implemented

------

## status
Show manager status [manager only]

### Manager state
Show service list on manager

```shell
cfy-go status state
```

### Manager version
Show manager version

```shell
cfy-go status version
```

------

## teardown
Teardown a manager [manager only]
* Not Implemented

------

## tenants
Handle Cloudify tenants (Premium feature)
* Not Implemented

------

## uninstall
Uninstall an application blueprint [manager only]
* Not Implemented

------

## user-groups
Handle Cloudify user groups (Premium feature)
* Not Implemented

------

## users
Handle Cloudify users
* Not Implemented

------

## workflows
Handle deployment workflows
* Not Implemented
