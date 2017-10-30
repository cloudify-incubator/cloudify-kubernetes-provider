ctx logger info "Install cfy-autoscale provider"
sudo cp /opt/cloudify-kubernetes-provider/src/k8s.io/autoscaler/cluster-autoscaler/cluster-autoscaler /usr/bin/cfy-autoscale
sudo chmod 555 /usr/bin/cfy-autoscale
sudo chown root:root /usr/bin/cfy-autoscale

ctx logger info "Create service"
sudo tee /etc/systemd/system/cfy-autoscale.service <<EOF
[Unit]
Description=cfy autoscale

[Service]
ExecStart=/usr/bin/cfy-autoscale --kubeconfig $HOME/.kube/config --cloud-config $HOME/cfy.json --cloud-provider cloudify --alsologtostderr
KillMode=process
Restart=on-failure
RestartSec=60s

[Install]
WantedBy=multi-user.target
EOF
sudo cp /etc/systemd/system/cfy-autoscale.service /etc/systemd/system/multi-user.target.wants/

ctx logger info "Start service"
sudo systemctl daemon-reload
sudo systemctl enable cfy-autoscale.service
sudo systemctl start cfy-autoscale.service

for retry_count in {1..10}
do
	status=`sudo systemctl status cfy-autoscale.service | grep "Active:"| awk '{print $2}'`
	ctx logger info "#${retry_count}: CFY Auto Scale state: ${status}"
	if [ "z$status" == 'zactive' ]; then
		break
	else
		ctx logger info "Wait little more."
		sleep 10
	fi
done
