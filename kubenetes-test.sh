# Based on https://github.com/kubernetes/community/blob/master/contributors/devel/development.md

mkdir -p kubenetes_checks
cd kubenetes_checks
export GOBIN=`pwd`/bin
export PATH=$PATH:`pwd`/bin
export GOPATH=`pwd`

go get -u github.com/tools/godep

# If your GOPATH has multiple paths, pick
# just one and use it instead of $GOPATH here.
# You must follow exactly this pattern,
# neither `$GOPATH/src/github.com/${your github profile name/`
# nor any other pattern will work.
working_dir=$GOPATH/src/k8s.io

user=kubernetes # fork profile

mkdir -p $working_dir
cd $working_dir

if [ ! -d "kubernetes" ]; then
	git clone https://github.com/$user/kubernetes.git
	# or: git clone git@github.com:$user/kubernetes.git

	cd $working_dir/kubernetes
	git remote add upstream https://github.com/kubernetes/kubernetes.git
	# or: git remote add upstream git@github.com:kubernetes/kubernetes.git

	# Never push to upstream master
	git remote set-url --push upstream no_push

	# Confirm that your remotes make sense:
	git remote -v
else
	cd kubernetes
	git pull
fi

cd $working_dir/kubernetes
make GOGCFLAGS="-N -l"

# no arm
# make cross

hack/install-etcd.sh

export PATH=${GOPATH}/src/k8s.io/kubernetes/third_party/etcd:${PATH}

cd $working_dir/kubernetes

# Run all the presubmission verification. Then, run a specific update script (hack/update-*.sh)
# for each failed verification. For example:
#   hack/update-gofmt.sh (to make sure all files are correctly formatted, usually needed when you add new files)
#   hack/update-bazel.sh (to update bazel build related files, usually needed when you add or remove imports)
make verify
