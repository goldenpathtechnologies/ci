#!/bin/bash

CI_INSTALL_DIR=~/.ci
CI_CMD=$CI_INSTALL_DIR/bin/ci

if [ -d "$CI_INSTALL_DIR" ]
then
  # Note: Version number extraction approach taken from https://stackoverflow.com/a/16817748/3141986
  CI_CURRENT_VERSION=$($CI_CMD -v | grep -Po '(?<=Version: )[\d\.]+')
  CI_NEW_VERSION=$(./bin/ci -v | grep -Po '(?<=Version: )[\d\.]+')

  # Note: version function taken from https://stackoverflow.com/a/37939589/3141986
  version() {
    echo "$@" | awk -F. '{ printf("%d%03d%03d%03d\n", $1,$2,$3,$4) }'
  }

  if [ "$(version "$CI_CURRENT_VERSION")" -ge "$(version "$CI_NEW_VERSION")" ]
  then
    echo "ci v$CI_CURRENT_VERSION is up to date"
    exit 0
  else
    ./uninstall.sh
  fi
fi

mkdir -p $CI_INSTALL_DIR/bin

cp ./bin/ci $CI_INSTALL_DIR/bin/ci
cp ./{LICENSE,CHANGELOG.md} $CI_INSTALL_DIR

{ echo "### BEGIN CI COMMAND";
cat ./ci.sh;
echo "### END CI COMMAND";} >> ~/.bashrc
