#!/bin/bash

# set -x
set -e

dir=$(cd $(dirname $0); pwd)
module=$(go list -m)

cd ${dir}/../ && CGO_ENABLED=0 go build -o ${module}

