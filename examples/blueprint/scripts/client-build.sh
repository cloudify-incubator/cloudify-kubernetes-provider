ctx logger info "Go to /opt"
cd /opt/

ctx logger info "Build client"
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`
go get github.com/0lvin-cfy/cloudify-rest-go-client/cfy-go
