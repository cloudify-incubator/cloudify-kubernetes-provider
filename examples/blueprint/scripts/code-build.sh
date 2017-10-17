ctx logger info "Build everything"
cd /opt/cloudify-kubernetes-provider

# kubernetes
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`

# cfy part
PACKAGEPATH=github.com/cloudify-incubator/cloudify-kubernetes-provider
ctx logger info "Build cfy-go"
sudo go install src/${PACKAGEPATH}/cfy-go/cfy-go.go
ctx logger info "Build cfy-kubernetes"
sudo go install src/cfy-kubernetes.go
