#!/usr/bin/env bash

# Using this as a quick workaround the missing exit 1 in gofmt when there's files with format errors
# $1 arg expected to be go path (e.g: /usr/local/go/bin/go)

gofmt_command="$1fmt"

gofmt_files=$(${gofmt_command} -l .)
if [[ -n ${gofmt_files} ]]; then
  echo 'The following files have format errors:'
  echo "${gofmt_files}"
  exit 1
fi

exit 0
