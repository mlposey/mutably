#!/bin/bash
# This script uses go test to find the total coverage percent for tested
# packages.
#
# E.g.,
# pkgA - 80% coverage
# pkgB - no tests
# pkgC - 85% coverage
# Total Coverage = (pkgA + pkgC) / 2 = 82.50

total_coverage="0"
package_count=0

# test_with_coverage recursively tests each subpackage of $1
test_with_coverage() {
  for item in $(ls $1); do
    if [[ -d "$1/$item" ]]; then
      test_with_coverage $1/$item
    fi
  done

  package_coverage=$(go test -coverprofile cp.out $1 | \
      grep -Po 'coverage: \K(\d{1,2}\.\d{1,2})')

  if [ -n "$package_coverage" ]; then
    package_count=$((package_count + 1))
    total_coverage=$(echo "$total_coverage + $package_coverage" | bc)
  fi
}

test_with_coverage .
echo "Total Coverage: $(echo "scale=2; $total_coverage / $package_count" | bc)"
