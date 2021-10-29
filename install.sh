#!/bin/bash
# TODO: Make this script idempotent
# TODO: Make version checks to ensure old versions do not overwrite new ones.
#  Old versions of the software can only be installed when the new version is
#  uninstalled.

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
