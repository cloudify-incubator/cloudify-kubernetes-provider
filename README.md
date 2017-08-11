# cloudify-rest-go-client

# install
```shell
sudo apt-get install gccgo-go golang-go
```

# run
```shell
go build src/status.go

./status blueprints list -host <your manager host> -user admin -password secret -tenant default_tenant
./status deployments list -host <your manager host> -user admin -password secret -tenant default_tenant
./status deployments create deployment -host <your manager host> -user admin -password secret -tenant default_tenant -blueprint blueprint
./status deployments delete  deployment -host <your manager host> -user admin -password secret -tenant default_tenant
./status executions list -host <your manager host> -user admin -password secret -tenant default_tenant
./status executions start uninstall -deployment deployment -host <your manager host> -user admin -password secret -tenant default_tenant
./status status state -host <your manager host> -user admin -password secret -tenant default_tenant
./status status version -host <your manager host> -user admin -password secret -tenant default_tenant
```
# reformat code
```shell
gofmt -w src/status.go
```
