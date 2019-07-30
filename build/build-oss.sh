#!/usr/bin/env bash
set -e

cd ${GOPATH}/src/minio/minio-csi/
GIT_SHA=`git rev-parse --short HEAD || echo "HEAD"`

rm -rf build/s3fs/csiplugin-connector build/s3fs/minio-csi

export GOARCH="amd64"
export GOOS="linux"

branch="v1.0.0"
version="v1.13.2"
commitId=$GIT_SHA
buildTime=`date "+%Y-%m-%d-%H:%M:%S"`

cd ${GOPATH}/src/minio/minio-csi/cmd
CGO_ENABLED=0 go build -ldflags "-X main._BRANCH_='$branch' -X main._VERSION_='$version-$commitId' -X main._BUILDTIME_='$buildTime'" -o minio-csi

cd ${GOPATH}/src/minio/minio-csi/build/s3fs/
CGO_ENABLED=0 go build csiplugin-connector.go

if [ "$1" == "" ]; then
  mv ${GOPATH}/src/minio/minio-csi/cmd/minio-csi ./
  docker build -t=k8s.gcr.io/csi-s3fsplugin:$version-$GIT_SHA ./
fi
