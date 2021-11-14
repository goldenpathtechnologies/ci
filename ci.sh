export CI_CMD=~/.ci/bin/ci

# TODO: Note that there is an existing 'ci' command:
#  http://manpages.ubuntu.com/manpages/precise/man1/rcsintro.1.html
#  Ensure there is a way to customize the name, or create a tool
#  that will update the bashrc with an alternate function name.
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
    # TODO: Note that 'return' returns the exit status of the last executed command: https://stackoverflow.com/a/43840545/3141986
    #  Remove the value from the return command.
    return 0
  else
    output=$($CI_CMD "$@")

    if [ -d "$output" ]
    then
      # TODO: A 'return' command with no parameters after the 'cd' command will suffice, remove the '|| return'.
      cd "$output" || return
      return 0
    else
      echo "$output"
      return 1
    fi
  fi
}
