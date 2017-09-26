#!/bin/bash

echo "mode: atomic" > coverage.out

govendor test -v -coverprofile=profile.out -covermode=atomic; __EXIT_CODE__=$?

if [ -f profile.out ]; then
  tail -n +2 profile.out >> coverage.out; rm profile.out
fi
