#!/bin/bash
set -e

mkdir -p builds

platforms=(
    "linux amd64"
    "darwin amd64"
    "darwin arm64"
    "windows amd64"
    "android arm64"
)

for platform in "${platforms[@]}"
do
    GOOS=$(echo $platform | cut -d' ' -f1)
    GOARCH=$(echo $platform | cut -d' ' -f2)
    output="synthera-$GOOS-$GOARCH"

    # Windows extension
    if [ "$GOOS" = "windows" ]; then
        output="$output.exe"
    fi

    echo "Building $output ..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o "builds/$output"
done
