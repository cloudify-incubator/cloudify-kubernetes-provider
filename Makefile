.PHONY: all
all: bin/cfy-go bin/cfy-kubernetes bin/cfy-autoscale

AUTOSCALEPACKAGE := k8s.io/autoscaler/cluster-autoscaler/cloudprovider
KUBERNETESPACKAGE := k8s.io/kubernetes/pkg/cloudprovider/providers/cloudifyprovider
PACKAGEPATH := github.com/cloudify-incubator/cloudify-rest-go-client

VERSION := `cd src/${PACKAGEPATH} && git rev-parse --short HEAD`

CLOUDPROVIDER ?= vsphere

.PHONY: reformat
reformat:
	rm -rfv pkg/*
	rm -rfv bin/*
	gofmt -w src/${PACKAGEPATH}/cloudify/rest/*.go
	gofmt -w src/${PACKAGEPATH}/cloudify/utils/*.go
	gofmt -w src/${PACKAGEPATH}/cloudify/tests/*.go
	gofmt -w src/${PACKAGEPATH}/cloudify/*.go
	gofmt -w src/${PACKAGEPATH}/cfy-go/*.go
	gofmt -w src/${PACKAGEPATH}/kubernetes/*.go
	# kubernetes parts
	gofmt -w src/${KUBERNETESPACKAGE}/*.go
	gofmt -w src/${AUTOSCALEPACKAGE}/cloudifyprovider/*.go
	gofmt -w src/*.go

define colorecho
	@tput setaf 2
	@echo -n $1
	@tput setaf 3
	@echo $2
	@tput sgr0
endef

# cloudify rest
CLOUDIFYREST := \
	src/${PACKAGEPATH}/cloudify/rest/rest.go \
	src/${PACKAGEPATH}/cloudify/rest/types.go

pkg/linux_amd64/${PACKAGEPATH}/cloudify/rest.a: ${CLOUDIFYREST}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudify/rest.a ${CLOUDIFYREST}

# cloudify kubernetes support
CLOUDIFYKUBERNETES := \
	src/${PACKAGEPATH}/kubernetes/mount.go \
	src/${PACKAGEPATH}/kubernetes/types.go

pkg/linux_amd64/${PACKAGEPATH}/kubernetes.a: ${CLOUDIFYKUBERNETES}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/kubernetes.a ${CLOUDIFYKUBERNETES}

# cloudify utils
CLOUDIFYUTILS := \
	src/${PACKAGEPATH}/cloudify/utils/utils.go

pkg/linux_amd64/${PACKAGEPATH}/cloudify/utils.a: ${CLOUDIFYUTILS}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudify/utils.a ${CLOUDIFYUTILS}

# cloudify
CLOUDIFYCOMMON := \
	src/${PACKAGEPATH}/cloudify/scalegroup.go \
	src/${PACKAGEPATH}/cloudify/scalenodes.go \
	src/${PACKAGEPATH}/cloudify/client.go \
	src/${PACKAGEPATH}/cloudify/agentfile.go \
	src/${PACKAGEPATH}/cloudify/nodes.go \
	src/${PACKAGEPATH}/cloudify/plugins.go \
	src/${PACKAGEPATH}/cloudify/instances.go \
	src/${PACKAGEPATH}/cloudify/loadbalancer.go \
	src/${PACKAGEPATH}/cloudify/events.go \
	src/${PACKAGEPATH}/cloudify/blueprints.go \
	src/${PACKAGEPATH}/cloudify/status.go \
	src/${PACKAGEPATH}/cloudify/executions.go \
	src/${PACKAGEPATH}/cloudify/deployments.go \
	src/${PACKAGEPATH}/cloudify/tenants.go

pkg/linux_amd64/${PACKAGEPATH}/cloudify.a: ${CLOUDIFYCOMMON} pkg/linux_amd64/${PACKAGEPATH}/cloudify/rest.a
	$(call colorecho,"Build: ",$@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudify.a ${CLOUDIFYCOMMON}

CFYGOLIBS := \
	pkg/linux_amd64/${PACKAGEPATH}/cloudify/utils.a \
	pkg/linux_amd64/${PACKAGEPATH}/kubernetes.a \
	pkg/linux_amd64/${PACKAGEPATH}/cloudify.a

# cfy-go
CFYGO := \
	src/${PACKAGEPATH}/cfy-go/blueprints.go \
	src/${PACKAGEPATH}/cfy-go/deployments.go \
	src/${PACKAGEPATH}/cfy-go/events.go \
	src/${PACKAGEPATH}/cfy-go/executions.go \
	src/${PACKAGEPATH}/cfy-go/info.go \
	src/${PACKAGEPATH}/cfy-go/instances.go \
	src/${PACKAGEPATH}/cfy-go/kubernetes.go \
	src/${PACKAGEPATH}/cfy-go/main.go \
	src/${PACKAGEPATH}/cfy-go/nodes.go \
	src/${PACKAGEPATH}/cfy-go/plugins.go \
	src/${PACKAGEPATH}/cfy-go/scaling.go \
	src/${PACKAGEPATH}/cfy-go/tenants.go

bin/cfy-go: ${CFYGO} ${CFYGOLIBS}
	$(call colorecho,"Install: ", $@)
	# delete -s -w if you want to debug
	go build -v -ldflags "-s -w -X main.versionString=${VERSION}" -o bin/cfy-go ${CFYGO}

# cloudify provider
CLOUDIFYPROVIDER := \
	src/${KUBERNETESPACKAGE}/init.go \
	src/${KUBERNETESPACKAGE}/instances.go \
	src/${KUBERNETESPACKAGE}/loadbalancer.go \
	src/${KUBERNETESPACKAGE}/zones.go

pkg/linux_amd64/${KUBERNETESPACKAGE}.a: pkg/linux_amd64/${PACKAGEPATH}/cloudify.a ${CLOUDIFYPROVIDER}
	$(call colorecho,"Build: ",$@)
	go build -v -i -o pkg/linux_amd64/${KUBERNETESPACKAGE}.a ${CLOUDIFYPROVIDER}

bin/cfy-kubernetes: pkg/linux_amd64/${KUBERNETESPACKAGE}.a pkg/linux_amd64/${PACKAGEPATH}/cloudify.a src/cfy-kubernetes.go
	$(call colorecho,"Install: ", $@)
	# delete -s -w if you want to debug
	go install -v -ldflags "-s -w -X main.versionString=${VERSION}" src/cfy-kubernetes.go

CLUSTERAUTOSCALERPROVIDER := \
	src/${AUTOSCALEPACKAGE}/cloudifyprovider/init.go \
	src/${AUTOSCALEPACKAGE}/cloudifyprovider/node_group.go \
	src/${AUTOSCALEPACKAGE}/cloudifyprovider/scale_provider.go

pkg/linux_amd64/${AUTOSCALEPACKAGE}/cloudifyprovider.a: ${CLUSTERAUTOSCALERPROVIDER} pkg/linux_amd64/${PACKAGEPATH}/cloudify.a
	$(call colorecho,"Build: ",$@)
	go build -v -i -o pkg/linux_amd64/${AUTOSCALEPACKAGE}/cloudifyprovider.a ${CLUSTERAUTOSCALERPROVIDER}

CLUSTERAUTOSCALER := \
	src/k8s.io/autoscaler/cluster-autoscaler/main.go \
	src/k8s.io/autoscaler/cluster-autoscaler/version.go

bin/cfy-autoscale: pkg/linux_amd64/${PACKAGEPATH}/cloudify.a ${CLUSTERAUTOSCALER} pkg/linux_amd64/${AUTOSCALEPACKAGE}/cloudifyprovider.a
	$(call colorecho,"Install: ", $@)
	# delete -s -w if you want to debug
	go build -v -ldflags "-s -w -X main.ClusterAutoscalerVersion=${VERSION}" -o bin/cfy-autoscale ${CLUSTERAUTOSCALER}

upload:
	cfy blueprints upload -b kubernetes_cluster examples/cluster_blueprint/${CLOUDPROVIDER}.yaml

create-for-upload: all
	cp -v bin/cfy-kubernetes examples/cluster_blueprint/resources/cfy-kubernetes
	cp -v bin/cfy-autoscale examples/cluster_blueprint/resources/cfy-autoscale
	cp -v bin/cfy-go examples/cluster_blueprint/resources/cfy-go

.PHONY: test
test:
	go test -cover ./src/${PACKAGEPATH}/...
	go get github.com/golang/lint/golint
	golint ./src/${PACKAGEPATH}/...
	golint ./src/${KUBERNETESPACKAGE}/...
	golint ./src/cfy-kubernetes.go
	golint ./src/${AUTOSCALEPACKAGE}/cloudifyprovider/...

	cfy blueprint validate examples/cluster_blueprint/aws.yaml
	cfy blueprint validate examples/cluster_blueprint/azure.yaml
	cfy blueprint validate examples/cluster_blueprint/gcp.yaml
	cfy blueprint validate examples/cluster_blueprint/openstack.yaml
	cfy blueprint validate examples/cluster_blueprint/vsphere.yaml
