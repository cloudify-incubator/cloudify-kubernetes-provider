ctx logger info "Download cfy-go"
cp /opt/bin/cfy-go /usr/bin/cfy-go
sudo chmod 555 /usr/bin/cfy-go
sudo chown root:root /usr/bin/cfy-go

ctx logger info "Create cloudify mount script"
PLUGINDIR=/usr/libexec/kubernetes/kubelet-plugins/volume/exec/cloudify~mount/
sudo mkdir -p $PLUGINDIR

ctx logger info "Create cfy config"
sudo tee $PLUGINDIR/mount <<EOF
#!/bin/bash
echo \$@ >> /var/log/mount-calls.log
/usr/bin/cfy-go kubernetes \$1 \$2 \$3 -deployment "$(ctx deployment id)" -instance "$(ctx instance id)" -tenant "${CFY_TENANT}" -password "${CFY_PASSWORD}" -user "${CFY_USER}" -host "${CFY_HOST}"
EOF

sudo chmod 555 -R $PLUGINDIR
sudo chown root:root -R $PLUGINDIR
