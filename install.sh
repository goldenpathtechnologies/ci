#!/bin/bash
# TODO: Make this script idempotent
# Step 1: Download files to /usr/local/ci
# - Create /usr/local/ci directory
# - Download release into /usr/local/ci/bin
# - Optional: Download source to /usr/local/ci/src
mkdir ~/.ci

# TODO: Use wget or curl to download release from GitHub
cp ./bin/go_build_ci_linux ~/.ci

# Step 2: Update ~/.bashrc with ci function
# - export envvar with path to ci executable
# - add function
# - run `source ~/.bashrc`

# TODO: Use ci.sh file from tagged release on remote
{ echo "### BEGIN CI COMMAND";
cat ./ci.sh;
echo "### END CI COMMAND";
echo ""; } >> ~/.bashrc

# shellcheck source=/dev/null
source ~/.bashrc
