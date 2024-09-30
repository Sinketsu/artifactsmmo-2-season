#!/bin/bash

set -e

ssh 51.250.80.190 'docker logs -n 10 -f $(docker ps -aqf "name=artifactsmmo")'
