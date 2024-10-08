#!/bin/bash

set -e

uuid=$(uuidgen)

docker build -t cr.yandex/crpa758jtdejs3rrqb7u/artifactsmmo:${uuid} -f cmd/player/Dockerfile .
docker push cr.yandex/crpa758jtdejs3rrqb7u/artifactsmmo:${uuid}
docker image rm cr.yandex/crpa758jtdejs3rrqb7u/artifactsmmo:${uuid}

yc compute instance update-container fhm4jlf2a9ovph7vfsmp  \
    --container-image=cr.yandex/crpa758jtdejs3rrqb7u/artifactsmmo:${uuid}
