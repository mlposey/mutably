#!/bin/bash

# test_with_coverage recursively tests each subpackage of $1
test_with_coverage() {
  for item in $(ls $1); do
    if [[ -d "$1/$item" ]]; then
      test_with_coverage $1/$item
    fi
  done

  go test -coverprofile cp.out $1
}

test_with_coverage .
