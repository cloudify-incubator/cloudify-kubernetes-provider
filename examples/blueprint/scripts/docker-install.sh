# before run: %wheel    ALL=(ALL)   NOPASSWD: ALL
# https://docs.docker.com/engine/installation/linux/centos/
# little cleanup
ctx logger info "Update basic instance"

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
sudo systemctl enable docker.service
sudo systemctl start docker
# reload user
exit
