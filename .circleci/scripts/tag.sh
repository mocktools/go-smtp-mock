#!/bin/sh
set -e

semver_regex_pattern="[0-9]+\.[0-9]+\.[0-9]+"

latest_changelog_tag() {
  grep -Po "(?<=\#\# \[)$semver_regex_pattern?(?=\])" CHANGELOG.md | head -n 1
}

latest_git_tag() {
  git tag --sort=v:refname | grep -E "v$semver_regex_pattern" | tail -n 1
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
  echo "Updating develop branch with new semver tag..."
  git checkout develop
  git merge "$tag_candidate" --ff --no-edit
  git push origin develop
else
  echo "Latest changelog tag ($tag_candidate) already released on GitHub. Tagging is not required."
fi
