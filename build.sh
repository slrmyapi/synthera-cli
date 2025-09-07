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

arch_map() {
    case "$1" in
        amd64) echo "x86_64" ;;
        386) echo "i386" ;;
        arm64) echo "aarch64" ;;
        arm) echo "arm" ;;
        *) echo "$1" ;;  # fallback
    esac
}

for platform in "${platforms[@]}"
do
    GOOS=$(echo $platform | cut -d' ' -f1)
    GOARCH=$(echo $platform | cut -d' ' -f2)
    uname_arch=$(arch_map $GOARCH)

    output="synthera-$GOOS-$uname_arch"

    if [ "$GOOS" = "windows" ]; then
        output="$output.exe"
    fi

    echo "Building $output ..."

    if [ "$GOOS" = "android" ]; then
        GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=1 CC=/opt/android-ndk/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang go build -ldflags="-s -w" -o "builds/$output"
    else
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "builds/$output"
    fi
done
