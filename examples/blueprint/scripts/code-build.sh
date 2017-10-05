ctx logger info "Build everything"
cd /opt/cloudify-rest-go-client

# kubernetes
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`

# cfy part
PACKAGEPATH=github.com/0lvin-cfy/cloudify-rest-go-client
ctx logger info "Build cfy-go"
go install src/${PACKAGEPATH}/cfy-go/cfy-go.go
ctx logger info "Build cfy-kubernetes"
go install src/cfy-kubernetes.go
ctx logger info "Build cfy-mount"
go install src/${PACKAGEPATH}/cfy-mount/cfy-mount.go
