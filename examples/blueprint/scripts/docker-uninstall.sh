sudo systemctl stop docker || ctx logger info "You dont have docker? wait several moments"
sudo yum remove -y -q docker-engine || ctx logger info "No docker yet"
sudo rm -f /etc/yum.repos.d/docker.repo
