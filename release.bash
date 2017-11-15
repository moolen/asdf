#!/bin/bash
set -xe

NEXT_VERSION=$(./asdf next-version)

# BRANCH CONFIGURATION
branch=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
devBranch=develop
MASTER=master
RELEASE=release-$NEXT_VERSION

git checkout -b $RELEASE $devBranch

# generate next version and changelog
./asdf generate
 
# commit version number increment
git commit -am "version bump to $NEXT_VERSION"

# merge release into master
git checkout $MASTER
git merge --no-ff $RELEASE
 
# create tag for new version from -master
git tag $NEXT_VERSION
 
# merge release branch back into develop
git checkout $devBranch
git merge --no-ff $RELEASE
 
# remove release branch
git branch -d $RELEASE