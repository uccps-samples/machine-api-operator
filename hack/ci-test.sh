#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE}")/..

OPENSHIFT_CI=${OPENSHIFT_CI:-""}
ARTIFACT_DIR=${ARTIFACT_DIR:-""}

function runTestsCI() {
  echo "CI env detected, run tests with jUnit report extraction"
  if [ -n "$ARTIFACT_DIR" ] && [ -d "$ARTIFACT_DIR" ]; then
    JUNIT_LOCATION="$ARTIFACT_DIR"/junit_machine_api_operator.xml
    echo "jUnit location: $JUNIT_LOCATION"
    go install -mod= github.com/jstemmer/go-junit-report@latest
    go test -v ./pkg/... ./cmd/... | tee >(go-junit-report > "$JUNIT_LOCATION")
  else
    echo "\$ARTIFACT_DIR not set or does not exists, no jUnit will be published"
    make unit
  fi
}


cd $REPO_ROOT && \
  source ./hack/fetch_ext_bins.sh && \
  fetch_tools && \
  setup_envs && \
if [ "$OPENSHIFT_CI" == "true" ]; then # detect ci environment there
  runTestsCI
else
  make unit
fi
