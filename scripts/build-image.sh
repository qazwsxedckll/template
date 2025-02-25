#!/bin/bash

# set -x
set -e

Help()
{
  # Display Help
  echo "Usage: ./build-image.sh <dockerfile path> [TAG]"
  echo "example:./build-image.sh ./Dockerfile latest"
}

[ ! -f $1 ] && Help && exit 1

dir="$( cd "$( dirname "$0"  )" && pwd  )"
workspace=${dir}/../
module=$(go list -m)

case $# in
  1 )
    dockerfile=$1; tag=latest
    ;;
  2 )
    dockerfile=$1; tag=$2
    ;;
  * ) Help && exit 1
esac

full_name=${module}:${tag}

docker build -t ${full_name} -f ${dockerfile} ${workspace}