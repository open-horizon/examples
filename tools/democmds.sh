#!/bin/bash

# Facilitate a CLI-based demo by reading a list of commands and displaying and running them 1 at a time

if [[ -z "$1" ]]; then
    cat << EndOfMessage
Usage: ${0##*/} <cmd-list-file>"

Will first show a \$ prompt. Press enter and it will show the first cmd in the file. Press enter and it will run that cmd. And so on....

Arguments:
  <cmd-list-file>  A list of cmds that should be displayed and run. Blank lines and lines that start with # will be ignored. Note: it is important for the final newline to be by itself (not on the line with the last command).
EndOfMessage
    exit
fi

cmdFile="$1"

if [[ ! -f "$cmdFile" ]]; then
    echo "Error: file $cmdFile does not exist."
    exit 1
fi

# Read in the cmds file
cmds=()
while IFS= read -r cmd; do
    pat1="^ *#"   # comments with #
    pat2="^ *$"   # blanks lines, or only spaces
    if [[ "$cmd" =~ $pat1 || "$cmd" =~ $pat2 ]]; then continue; fi
    cmds+=("$cmd")
done < "$cmdFile"
if [[ ${#cmds[@]} -eq 0 ]]; then
    echo "Error: only empty or commented out lines in $cmdFile"
    exit 2
fi

# Display and execute 1 cmd at a time
for cmd in "${cmds[@]}"; do
    printf '\e[0;32m$ \e[0m'   # display prompt so we wait before displaying the next cmd
    read junk   # wait for them to hit enter
    printf '\e[0;32m$ \e[0m'"$cmd"   # display the cmd with no newline
    read junk   # wait for them to hit enter
    eval "$cmd"   # eval the cmd so it runs in this process, in case it is setting a variable
done
