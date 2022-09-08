#!/bin/sh
set -e

RELEASES_URL="https://github.com/mocktools/go-smtp-mock/releases"
GO_PKG_URL="https://pkg.go.dev/fetch/github.com/mocktools/go-smtp-mock"
SUCCESS_MESSAGE="$(tput bold)$(tput setaf 2)[SUCCESS]$(tput sgr0) Latest smtpmock release has been publishid on Go reference"

latest_release() {
  curl -sL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" | rev | cut -f1 -d'/'| rev
}

publish_release() {
  if [[ $(curl -s -X POST "$GO_PKG_URL@$(latest_release)") =~ "could not be found" ]]; then exit 1
  else echo $SUCCESS_MESSAGE
  fi
}

publish_release
