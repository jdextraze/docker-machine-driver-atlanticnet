#!/bin/bash

RELEASE_TAG=$1
RELEASE_NAME=$2
RELEASE_DESCRIPTION=$3

gox -os="darwin linux windows" -arch="amd64" -output "bin/{{.OS}}/{{.Arch}}/{{.Dir}}"

mkdir -p releases

tar -cjf releases/darwin-amd64-release.tar.bz2 bin/darwin/amd64/docker-machine-driver-atlanticnet
tar -cjf releases/linux-amd64-release.tar.bz2 bin/linux/amd64/docker-machine-driver-atlanticnet
zip releases/windows-amd64-release.zip bin/windows/amd64/docker-machine-driver-atlanticnet.exe

github-release release \
    --user jdextraze \
    --repo docker-machine-driver-atlanticnet \
    --tag "${RELEASE_TAG}" \
    --name "${RELEASE_NAME}" \
    --description "${RELEASE_DESCRIPTION}" \
    --pre-release

github-release upload \
    --user jdextraze \
    --repo docker-machine-driver-atlanticnet \
    --tag ${RELEASE_TAG} \
    --name "docker-machine-driver-atlanticnet-darwin-amd64-${RELEASE_TAG}.tar.bz2" \
    --file releases/darwin-amd64-release.tar.bz2

github-release upload \
    --user jdextraze \
    --repo docker-machine-driver-atlanticnet \
    --tag ${RELEASE_TAG} \
    --name "docker-machine-driver-atlanticnet-linux-amd64-${RELEASE_TAG}.tar.bz2" \
    --file releases/linux-amd64-release.tar.bz2

github-release upload \
    --user jdextraze \
    --repo docker-machine-driver-atlanticnet \
    --tag ${RELEASE_TAG} \
    --name "docker-machine-driver-atlanticnet-windows-amd64-${RELEASE_TAG}.zip" \
    --file releases/windows-amd64-release.zip

rm -rf bin
rm -rf releases
