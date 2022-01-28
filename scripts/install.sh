#!/bin/bash

ORIG_DIR=$(pwd)
INSTALL_DIR=$ORIG_DIR

if [[ $ORIG_DIR == */scripts ]]
then
  # Note: Simple way to get the parent directory, https://stackoverflow.com/a/42956288/3141986
  INSTALL_DIR=$(builtin cd .. && pwd)
fi

if [ ! -f "$INSTALL_DIR/bin/ci" ]
then
  echo "ci executable not present, unable to install"
  exit 1
fi

cd "$INSTALL_DIR" || { echo "Install directory '$INSTALL_DIR' does not exist"; exit 1; }

CI_INSTALL_DIR=~/.ci
CI_CMD=$CI_INSTALL_DIR/bin/ci

if [ -d "$CI_INSTALL_DIR" ]
then
  # Note: Version number extraction approach taken from https://stackoverflow.com/a/16817748/3141986
  CI_CURRENT_VERSION=$($CI_CMD -v | grep -Po '(?<=Version: )[\d\.]+.*')
  CI_NEW_VERSION=$(./bin/ci -v | grep -Po '(?<=Version: )[\d\.]+.*')

  # Note: version function taken from https://stackoverflow.com/a/37939589/3141986
  version() {
    echo "$@" | awk -F. '{ printf("%d%03d%03d%03d\n", $1,$2,$3,$4) }'
  }

  if [ "$(version "$CI_CURRENT_VERSION")" -gt "$(version "$CI_NEW_VERSION")" ]
  then
    echo "ci is already installed at v$CI_CURRENT_VERSION."
    echo "Please uninstall the current version before installing v$CI_NEW_VERSION."
    exit 0
  elif [ "$(version "$CI_CURRENT_VERSION")" -eq "$(version "$CI_NEW_VERSION")" ]
  then
    echo "The installed version of ci ($CI_CURRENT_VERSION) is up to date."
    exit 0
  else
    ./uninstall.sh
  fi
fi

mkdir -p $CI_INSTALL_DIR/bin

cp ./bin/ci $CI_INSTALL_DIR/bin/ci
cp ./{LICENSE,CHANGELOG.md} $CI_INSTALL_DIR

{ echo "### BEGIN CI COMMAND";
cat ./scripts/ci.sh;
echo "### END CI COMMAND";} >> ~/.bashrc

cd "$ORIG_DIR" || exit