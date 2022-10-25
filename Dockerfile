FROM openshift/golang-builder@sha256:4820580c3368f320581eb9e32cf97aeec179a86c5749753a14ed76410a293d83 AS builder
ENV __doozer=update BUILD_RELEASE=202202160023.p0.g668c5b5.assembly.stream BUILD_VERSION=v4.10.0 OS_GIT_MAJOR=4 OS_GIT_MINOR=10 OS_GIT_PATCH=0 OS_GIT_TREE_STATE=clean OS_GIT_VERSION=4.10.0-202202160023.p0.g668c5b5.assembly.stream SOURCE_GIT_TREE_STATE=clean 
ENV __doozer=merge OS_GIT_COMMIT=668c5b5 OS_GIT_VERSION=4.10.0-202202160023.p0.g668c5b5.assembly.stream-668c5b5 SOURCE_DATE_EPOCH=1643021182 SOURCE_GIT_COMMIT=668c5b52b104f44f339622cd397c51f584879e02 SOURCE_GIT_TAG=v0.2.0-1424-g668c5b52 SOURCE_GIT_URL=https://github.com/uccps-samples/machine-api-operator 
WORKDIR /go/src/github.com/uccps-samples/machine-api-operator
COPY . .
RUN NO_DOCKER=1 make build

FROM openshift/ose-base:v4.10.0.20220216.010142
ENV __doozer=update BUILD_RELEASE=202202160023.p0.g668c5b5.assembly.stream BUILD_VERSION=v4.10.0 OS_GIT_MAJOR=4 OS_GIT_MINOR=10 OS_GIT_PATCH=0 OS_GIT_TREE_STATE=clean OS_GIT_VERSION=4.10.0-202202160023.p0.g668c5b5.assembly.stream SOURCE_GIT_TREE_STATE=clean 
ENV __doozer=merge OS_GIT_COMMIT=668c5b5 OS_GIT_VERSION=4.10.0-202202160023.p0.g668c5b5.assembly.stream-668c5b5 SOURCE_DATE_EPOCH=1643021182 SOURCE_GIT_COMMIT=668c5b52b104f44f339622cd397c51f584879e02 SOURCE_GIT_TAG=v0.2.0-1424-g668c5b52 SOURCE_GIT_URL=https://github.com/uccps-samples/machine-api-operator 
COPY --from=builder /go/src/github.com/uccps-samples/machine-api-operator/install manifests
COPY --from=builder /go/src/github.com/uccps-samples/machine-api-operator/bin/machine-api-operator .
COPY --from=builder /go/src/github.com/uccps-samples/machine-api-operator/bin/nodelink-controller .
COPY --from=builder /go/src/github.com/uccps-samples/machine-api-operator/bin/machine-healthcheck .
COPY --from=builder /go/src/github.com/uccps-samples/machine-api-operator/bin/machineset ./machineset-controller
COPY --from=builder /go/src/github.com/uccps-samples/machine-api-operator/bin/vsphere ./machine-controller-manager

LABEL \
        io.openshift.release.operator="true" \
        name="openshift/ose-machine-api-operator" \
        com.redhat.component="ose-machine-api-operator-container" \
        io.openshift.maintainer.product="OpenShift Container Platform" \
        io.openshift.maintainer.component="Cloud Compute" \
        release="202202160023.p0.g668c5b5.assembly.stream" \
        io.openshift.build.commit.id="668c5b52b104f44f339622cd397c51f584879e02" \
        io.openshift.build.source-location="https://github.com/uccps-samples/machine-api-operator" \
        io.openshift.build.commit.url="https://github.com/uccps-samples/machine-api-operator/commit/668c5b52b104f44f339622cd397c51f584879e02" \
        version="v4.10.0"

