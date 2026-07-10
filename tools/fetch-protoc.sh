#!/bin/sh

case "$(arch)" in
  "x86_64") export ARCH="x86_64"
  ;;
  "aarch64") export ARCH="aarch_64"
  ;;
  "arm64") export ARCH="aarch_64"
  ;;
esac

PROTOC_VERSION='31.1'

export SOURCE="https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-linux-$ARCH.zip"

echo $SOURCE

curl -L $SOURCE > /tmp/protoc.zip

unzip /tmp/protoc.zip -d $HOME/local
