export CI_CMD=~/.ci/bin/ci

# TODO: Note that there is an existing 'ci' command:
#  http://manpages.ubuntu.com/manpages/precise/man1/rcsintro.1.html
#  Ensure there is a way to customize the name, or create a tool
#  that will update the bashrc with an alternate function name.
#  RCS is obscure/old enough that it isn't worth it to implement
#  this improvement until enough people complain about it.
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
    return
  else
    output=$($CI_CMD "$@")
    lastCode=$?

    if [ -d "$output" ]
    then
      # shellcheck disable=SC2164
      cd "$output"
      return
    elif [ "$lastCode" == 0 ] && [ -z "$output" ]
    then
      echo "Program forcefully exited"
      return "$lastCode"
    else
      echo "$output"
      return "$lastCode"
    fi
  fi
}
