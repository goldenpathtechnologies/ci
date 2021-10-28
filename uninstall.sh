#!/bin/bash
# TODO: Make this script idempotent
# Step 1: Remove entry from ~/.bashrc
# - use sed or any tool that is builtin to remove function and export command
# - unset variable
# - run `source ~/.bashrc`
sed -i '/### BEGIN CI COMMAND/,/### END CI COMMAND/d' ~/.bashrc

unset CI_CMD

# shellcheck source=/dev/null
source ~/.bashrc

# Step 2: Remove ~/.ci
# - run `rm -rf ~/.ci`
rm -rf ~/.ci
