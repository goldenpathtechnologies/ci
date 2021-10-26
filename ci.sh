#!/bin/bash

ci() {
  exitArgs=("-v" "--version" "-h" "--help")
  containsExitArgs=false

  for arg in "$@"
  do
    for i in "${exitArgs[@]}"
    do
      if [ "$arg" == "$i" ]
      then
        containsExitArgs=true
      fi
    done
  done

  if [ "$containsExitArgs" = true ]
  then
    ./bin/go_build_ci_linux "$@"
    return 0
  else
    output=$(./bin/go_build_ci_linux "$@")

    if [ -d "$output" ]
    then
      cd "$output" || return
      return 0
    else
      echo "$output"
      return 1
    fi
  fi
}