# set -x
set -e

dir="$( cd "$( dirname "$0"  )" && pwd  )"
workspace=${dir}/../
module=$(go list -m)

cd ${workspace}
mkdir -p release/
rm -rf release/*

mkdir -p release/${module}
cp --parent configs/config.example.toml release/${module}
cp README.md release/${module}

go test ./...
version=$(git describe --tags --abbrev=0 --always)
go build -ldflags="-X ${module}/cmd.Version=${version}" -o release/${module}

tar -czvf release/${module}-${version}.tar.gz -C release ${module}/
