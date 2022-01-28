#!/bin/bash

CI_INSTALL_DIR=~/.ci

if [ ! -d "$CI_INSTALL_DIR" ]
then
  echo "ci is not installed"
  exit 0
fi

sed -i '/### BEGIN CI COMMAND/,/### END CI COMMAND/d' ~/.bashrc

unset CI_CMD

rm -rf ~/.ci
