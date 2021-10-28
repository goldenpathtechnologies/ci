export CI_CMD="./bin/go_build_ci_linux"

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
    $CI_CMD "$@"
    return 0
  else
    output=$($CI_CMD "$@")

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