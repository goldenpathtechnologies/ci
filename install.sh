#!/bin/bash

# Step 1: Download files to /usr/local/ci
# - Create /usr/local/ci directory
# - Download release into /usr/local/ci/bin
# - Optional: Download source to /usr/local/ci/src
mkdir /usr/local/ci
mkdir /usr/local/ci/bin

# TODO: Use wget or curl to download release from GitHub
cp ./bin/go_build_ci_linux /usr/local/ci/bin

# Step 2: Update ~/.bashrc with ci function
# - export envvar with path to ci executable
# - add function
# - run `source ~/.bashrc`

# TODO: Use ci.sh.tpl file from tagged release
{ echo "### Begin Ci function";
cat ./ci.sh.tpl;
echo "### End Ci function"; } >> ~/.bashrc

# shellcheck source=/dev/null
source ~/.bashrc
