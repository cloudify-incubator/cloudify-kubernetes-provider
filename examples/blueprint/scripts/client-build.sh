ctx logger info "Go to /opt"
cd /opt/

ctx logger info "Build client"
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`
go get github.com/cloudify-incubator/cloudify-rest-go-client/cfy-go
