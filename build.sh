#!/bin/bash

#TODO:begin checking before build
go_module=ihttp-example
app_name=ihttp-example
http_port=6001
https_port=6443

github_repo="ghcr.io"
github_token="xxxxx"
github_image_name="xxxx/xxxxx"
github_user="xxxx"
#TODO:end checking before build

#Read argument
for ARGUMENT in "$@"; do
  KEY=$(echo $ARGUMENT | cut -f1 -d=)

  KEY_LENGTH=${#KEY}
  VALUE="${ARGUMENT:$KEY_LENGTH+1}"

  export "$KEY"="$VALUE"
done

#Version checking
if [ -z "$version" ]; then
  echo "Usage: ./build.sh version=BUILD_VERSION"
  exit 1
fi

github_repo_image_name_with_version_tag="$github_repo/$github_image_name:$version"

echo "Build Image:"
echo "GitHub Image --> $github_repo_image_name_with_version_tag"


#Login to all repository
docker login $github_repo -u $github_user -p $github_token

#build process
rm -rf vendor
go env -w GOSUMDB=off
rm go.mod
rm go.sum
go mod init $go_module
go mod tidy
go mod vendor

docker buildx build \
--build-arg APP_NAME=$app_name \
--build-arg HTTP_PORT=$http_port \
--build-arg HTTPS_PORT=$https_port \
--build-arg version=$version  \
--platform=linux/amd64,linux/arm64 \
--output=type=image,push=true \
--no-cache \
--tag "$github_repo_image_name_with_version_tag"  .

echo "last build --> $version" >> lastbuild.txt
