#!/bin/bash

APP_NAME=meshblu-device-claimer
TMP_DIR=$PWD/tmp/cross
IMAGE_NAME=local/$APP_NAME

build() {
  for goos in darwin linux windows; do
    for goarch in 386 amd64; do
      local filename="${APP_NAME}-${goos}-${goarch}"
      echo "building: ${filename}"
      env GOOS="$goos" GOARCH="$goarch" go build -a -ldflags '-s' -o "${TMP_DIR}/${filename}"
    done
  done
}

init() {
  rm -rf $TMP_DIR/ \
   && mkdir -p $TMP_DIR/
}

panic() {
  local message=$1
  echo $message
  exit 1
}

main() {
  init    || panic "init failed"
  build   || panic "build failed"
}
main
