#!/bin/sh
set -e

semver_regex_pattern="[0-9]+\.[0-9]+\.[0-9]+"

latest_changelog_tag() {
  grep -Po "(?<=\#\# \[)$semver_regex_pattern?(?=\])" CHANGELOG.md | cut -d"-" -f 1 | head -n 1
}

latest_git_tag() {
  git tag -l | grep -E "^v$semver_regex_pattern" | cut -d"-" -f 1 | sort | tail -n 1
}

tag_candidate="v$(latest_changelog_tag)"

if [ "$tag_candidate" != "$(latest_git_tag)" ]
then
  echo "Configuring git..."
  git config --global user.email "${PUBLISHER_EMAIL}"
  git config --global user.name "${PUBLISHER_NAME}"
  echo "Pushing new semver tag to GitHub..."
  git tag "$tag_candidate"
  git push --tags
else
  echo "Latest changelog tag ($tag_candidate) already released on GitHub. Tagging is not required."
fi
