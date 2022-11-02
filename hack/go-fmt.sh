#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  for TARGET in "${@}"; do
    find "${TARGET}" -name '*.go' ! -path '*/vendor/*' ! -path '*/.build/*' -exec gofmt -s -w {} \+
  done
  git diff --exit-code
else
  docker run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/uccps-samples/machine-api-operator:z" \
    --workdir /go/src/github.com/uccps-samples/machine-api-operator \
    --env GO111MODULE="$GO111MODULE" \
    --env GOFLAGS="$GOFLAGS" \
    openshift/origin-release:golang-1.13 \
    ./hack/go-fmt.sh "${@}"
fi
