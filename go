#!/bin/bash

if grep -q '^/tmp/' <<<$PWD; then
  rsync -ruhv --exclude="*.go" $REPO_DIR/ $PWD > /dev/null
fi

$GO_ORIG_BIN "$@"
