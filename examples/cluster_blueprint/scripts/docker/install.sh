#!/bin/bash

# before run: %wheel    ALL=(ALL)   NOPASSWD: ALL
# https://docs.docker.com/engine/installation/linux/centos/
# little cleanup
ctx logger info "Update basic instance"
VM_VERSION=`grep -w '^NAME=' /etc/os-release`

if [[ "$VM_VERSION" == 'NAME="CentOS Linux"' ]]; then
	sudo yum install deltarpm epel-release unzip -q -y
	sudo yum update -y -q

	ctx logger info "Enable docker"

	# enable docker
	sudo tee /etc/yum.repos.d/docker.repo <<-'EOF'
	[dockerrepo]
	name=Docker Repository
	baseurl=https://yum.dockerproject.org/repo/main/centos/7/
	enabled=1
	gpgcheck=1
	gpgkey=https://yum.dockerproject.org/gpg
	EOF

	# add users
	sudo groupadd docker || ctx logger info "Docker group already exist?"
	sudo usermod -aG docker centos  || ctx logger info "User already in docker group?"

	# install docker
	ctx logger info "Update repos"
	sudo yum update -y -q
	ctx logger info "Install docker"
	sudo yum install docker-engine-1.12.6 -y -q
elif [[ "$VM_VERSION" == 'NAME="Ubuntu"' ]]; then
	apt-get update && apt-get install -y apt-transport-https curl
	curl -s https://download.docker.com/linux/ubuntu/gpg | apt-key add -
	sudo add-apt-repository \
	   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
	   $(lsb_release -cs) \
	   stable"
	apt-get update
	sudo apt-get install -y docker.io
	sudo docker run hello-world
else
	ctx logger info "Unknow OS"
fi

sudo systemctl enable docker.service
sudo systemctl start docker
# reload user
exit
