#!/bin/sh
set -e

RELEASES_URL="https://github.com/mocktools/go-smtp-mock/releases"
BINARY_NAME="smtpmock"
ARCH_TYPE=".tar.gz"
TAR_FILE="$BINARY_NAME$ARCH_TYPE"

latest_release() {
  curl -sL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" | rev | cut -f1 -d'/'| rev
}

remove_tmp_download() {
  rm -f "$TAR_FILE"
}

download() {
  test -z "$VERSION" && VERSION="$(latest_release)"
  test -z "$VERSION" && {
    echo "Unable to get smtpmock release." >&2
    exit 1
  }
  remove_tmp_download
  curl -s -L -o "$TAR_FILE" "$RELEASES_URL/download/$VERSION/smtpmock_$(uname -s)_$(uname -m)$ARCH_TYPE"
}

extract() {
  tar -zxf "$TAR_FILE" "$BINARY_NAME"
  remove_tmp_download
}

download
extract
