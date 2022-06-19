#!/bin/bash
set -e

# target folder
cd /target/utils

# run translations
go run .

# make sure we didn't change line endings by running
# this on docker on windows
find /target -type f -name '*.vdf' -exec unix2dos --quiet {} \;

# show the differences
git status

# debug
if [[ "${DEBUG}" != "" ]]; then
  tail -f /dev/null
fi