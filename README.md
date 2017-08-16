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
sudo apt-get install golang-go
```
# install

```shell
sudo apt-get install golang-go
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`
make all
```

# kubernetes
```shell
go get -d k8s.io/kubernetes
cd src/k8s.io/kubernetes
make
cd $GOPATH
```

# reformat code

```shell
make reformat
```
# Functionlity from original cfy client

* Common parameters:
    * `-host`: manager host
    * `-user`: manager user
    * `-password`: manager password
    * `-tenant`: manager tenant
* Example:

```shell
cfy-go status version -host <your manager host> -user admin -password secret -tenant default_tenant
```
* Limitations:
	* supported only simple version of inputs/outputs when we have map[string]string. (On API level, in cli inputs/outputs is not supported at all.)

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
* Not Implemented

### get
Retrieve blueprint information [manager only]
* Not Implemented

### inputs
Retrieve blueprint inputs [manager only]
* Not Implemented

### install-plugins
Install plugins [locally]
* Not Implemented

### list
List blueprints [manager only]
* Partially implemented, pagination is unsupported

```shell
cfy-go blueprints list
```

### package
Create a blueprint archive
* Not Implemented

### upload
Upload a blueprint [manager only]
* Not Implemented

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
* Partially implemented, set inputs is not supported. Internally supported only map string<->string.

```shell
cfy-go deployments create deployment  -blueprint blueprint
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
* Partially implemented, pagination is unsupported

```shell
cfy-go deployments list
```

### outputs
Show deployment outputs [manager only]
* Not Implemented

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
* Partially implemented, pagination is unsupported

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
* Partially implemented, pagination is unsupported

```shell
cfy-go executions list
cfy-go executions list -deployment deployment

```

### start
Execute a workflow [manager only]
* Partially implemented, set parametes is not supported.

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
* Not Implemented

------

## nodes
Handle a deployment's nodes
* Not Implemented

------

## plugins
Handle plugins on the manager
* Not Implemented

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
