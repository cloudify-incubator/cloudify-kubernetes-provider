#!/bin/bash

existing_docker_install=$(ctx instance runtime-properties existing_docker_install)

if [[ $existing_docker_install == 0 ]]; then
    ctx logger info "Docker Pre-installed."
    exit 0
fi

sudo systemctl stop docker || ctx logger info "You dont have docker? wait several moments"
sudo yum remove -y -q docker-engine || ctx logger info "No docker yet"
sudo rm -f /etc/yum.repos.d/docker.repo
