ctx logger info "Go to /opt"
cd /opt/

ctx logger info "Download top level sources"
# take ~ 16m34.350s for rebuild, 841M Disk Usage
rm -rf cloudify-rest-go-client
git clone https://github.com/cloudify-incubator/cloudify-kubernetes-provider.git --depth 1 -b master
sed -i "s|git@github.com:|https://github.com/|g" cloudify-kubernetes-provider/.gitmodules

cd cloudify-kubernetes-provider
ctx logger info "Download submodules sources"
git submodule init
git submodule update

ctx logger info "Update compiler"
sudo CGO_ENABLED=0 go install -a -installsuffix cgo std
