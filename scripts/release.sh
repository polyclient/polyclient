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

if [[ $# -ne 1 ]]; then
	usage
fi

VERSION=$1
is_semver "$VERSION"

TAG="v$VERSION"

echo "Creating tag: $TAG"
git tag -a "$TAG" -m "publish $TAG"
echo "✅ Tag $TAG created successfully"

echo "Pushing tag: $TAG"
git push origin "$TAG"
echo "✅ Tag $TAG pushed successfully"
