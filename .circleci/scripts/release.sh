#!/bin/sh
set -e

github_namespace=$1
github_repository=$2

get_latest_tag() {
  git tag -l | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+" | cut -d"-" -f 1 | sort | tail -n 1
}

latest_tag=$(get_latest_tag)
current_version="$(printf '%s' "$latest_tag" | cut -c 2-2)"
release_version=$(if [ "$current_version" -gt 1 ]; then echo "v$current_version"; else echo; fi)

publish_release() {
  echo "Triggering pkg.go.dev about new release..."
  curl -X POST "https://pkg.go.dev/fetch/github.com/$github_namespace/$github_repository/$release_version@$latest_tag"
}

publish_release
