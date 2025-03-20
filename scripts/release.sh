#!/usr/bin/env bash

set -euo pipefail

function usage() {
	echo "Usage: $0 <version>"
	echo "Example: $0 0.1.0"
	exit 1
}

function is_semver() {
	local version=$1
	if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]]; then
		echo "❌ Version must follow semantic versioning (e.g., 1.0.0, 1.0.0-alpha.1, etc.)"
		exit 1
	fi
}

function ensure_clean_git() {
	if ! git diff-index --quiet HEAD --; then
		echo "❌ Uncommitted changes detected. Commit or stash them before proceeding."
		exit 1
	fi
}

if [[ $# -ne 1 ]]; then
	usage
fi

VERSION=$1
is_semver "$VERSION"
ensure_clean_git

TAG="v$VERSION"

echo "Creating tag: $TAG"
git tag -a "$TAG" -m "publish $TAG"
echo "✅ Tag $TAG created successfully"

echo "Updating version.txt"
printf "%s" "$VERSION" | tee version.txt >/dev/null

echo "Pushing tag: $TAG"
git add version.txt
git commit -m "release: version $VERSION"
git push origin HEAD
git push origin "$TAG"
echo "✅ Tag $TAG pushed successfully"
