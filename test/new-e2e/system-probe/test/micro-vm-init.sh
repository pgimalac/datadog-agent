#!/bin/bash

set -eo xtrace
trap 'echo $SECONDS' DEBUG

GOVERSION=$1
RETRY_COUNT=$2
RUNNER_CMD="$(shift 2; echo "$*")"

# Add provisioning steps here !
## Set go version correctly
gimme "$GOVERSION"
# shellcheck source=/dev/null
source "$HOME/.gimme/envs/go$GOVERSION.env"
## Start docker
systemctl start docker
## Load docker images
if [ -f /docker-images.txt ]; then
  DOCKER_USERNAME=$(</docker-username)
  DOCKER_REGISTRY=$(</docker-registry)
  docker login --username "${DOCKER_USERNAME}" --password-stdin "${DOCKER_REGISTRY}" < /docker-password
  xargs -L1 -a /docker-images.txt docker pull
fi

# VM provisioning end !

# Start tests
rm -rf /ci-visibility
CODE=0
/test-runner -retry "$RETRY_COUNT" "$RUNNER_CMD" || CODE=$?

pushd /ci-visibility
tar czvf testjson.tar.gz testjson
tar czvf junit.tar.gz junit

exit $CODE
