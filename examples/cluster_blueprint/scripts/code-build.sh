ctx logger info "Build everything"
sudo mkdir -p /opt/cloudify-kubernetes-provider
cd /opt/cloudify-kubernetes-provider

# kubernetes
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`

# cfy part
PACKAGEPATH=github.com/cloudify-incubator/cloudify-kubernetes-provider

ctx logger info "Build cfy-go"
go install src/${PACKAGEPATH}/cfy-go/cfy-go.go

ctx logger info "Build cfy-kubernetes"
go install src/cfy-kubernetes.go

ctx logger info "Build cluster-autoscaler"
cd /opt/scaller/src/k8s.io/autoscaler/cluster-autoscaler/
export GOBIN=/opt/scaller/bin
export PATH=$PATH:/opt/scaller/bin
export GOPATH=/opt/scaller/
make
