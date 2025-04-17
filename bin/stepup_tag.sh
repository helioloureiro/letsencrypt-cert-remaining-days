#! /usr/bin/env bash
##
##
last_tag=$(git tag -l --sort taggerdate | tail -1)
major=$(echo "$last_tag" | cut -d. -f1)
minor=$(echo "$last_tag" | cut -d. -f2)
patchlevel=$(echo "$last_tag" | cut -d. -f3)

patchlevel=$(($patchlevel+1))

echo "Adding tag: $major.$minor.$patchlevel"
git tag "$major.$minor.$patchlevel"