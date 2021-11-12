#!/bin/bash

CI_INSTALL_DIR=~/.ci

if [ ! -d "$CI_INSTALL_DIR" ]
then
  echo "ci is not installed"
  exit 0
fi

# TODO: Ensure that newlines created by the installation script are removed from the .bashrc here
sed -i '/### BEGIN CI COMMAND/,/### END CI COMMAND/d' ~/.bashrc

unset CI_CMD

# shellcheck source=/dev/null
source ~/.bashrc

rm -rf ~/.ci
