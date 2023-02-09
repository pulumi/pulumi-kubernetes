#!/bin/bash
# Usage ./install [NAME] [VERSION] [REPO] [BUNDLE]
NAME=$1
VERSION=$2
REPO=$3
BUNDLE=$4
PLAT=$(go env GOOS)-$(go env GOARCH)
BIN=$NAME$(go env GOEXE)
URL="https://github.com/pulumi/$REPO/releases/download/v$VERSION/$BUNDLE-v$VERSION-$PLAT.tar.gz"
echo "Installing $BIN $VERSION from $URL"
wget -q -O - "$URL" | tar -xzf - "$BIN"
