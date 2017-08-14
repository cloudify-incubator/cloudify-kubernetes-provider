# cloudify-rest-go-client

# install
```shell
sudo apt-get install gccgo-go golang-go
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`
make all
```

# run
```shell
cfy-go blueprints list -host <your manager host> -user admin -password secret -tenant default_tenant
cfy-go deployments list -host <your manager host> -user admin -password secret -tenant default_tenant
cfy-go deployments create deployment -host <your manager host> -user admin -password secret -tenant default_tenant -blueprint blueprint
cfy-go deployments delete  deployment -host <your manager host> -user admin -password secret -tenant default_tenant
cfy-go executions list -host <your manager host> -user admin -password secret -tenant default_tenant
cfy-go executions start uninstall -deployment deployment -host <your manager host> -user admin -password secret -tenant default_tenant
cfy-go status state -host <your manager host> -user admin -password secret -tenant default_tenant
cfy-go status version -host <your manager host> -user admin -password secret -tenant default_tenant
```
# reformat code
```shell
make reformat
```
