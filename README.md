# cloudify-rest-go-client

# install
```shell
sudo apt-get install gccgo-go golang-go
```

# run
```shell
go run src/status.go -host <your manager host> -user admin -password secret -tenant default_tenant -command status
```
# reformat code
```shell
gofmt -w src/status.go
```
