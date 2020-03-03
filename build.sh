#!/bin/bash

# Remove old binaries (if any)
rm -rf dist

env GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o "dist/funstream_linux_i386" ./cmd/funstream/funstream.go             # Linux i386
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "dist/funstream_linux_x86_64" ./cmd/funstream/funstream.go         # Linux 64bit
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o "dist/funstream_linux_arm" ./cmd/funstream/funstream.go      # Linux armv5/armel/arm (it also works on armv6)
env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o "dist/funstream_linux_armhf" ./cmd/funstream/funstream.go    # Linux armv7/armhf
env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o "dist/funstream_linux_aarch64" ./cmd/funstream/funstream.go        # Linux armv8/aarch64
env GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o "dist/funstream_freebsd_x86_64" ./cmd/funstream/funstream.go     # FreeBSD 64bit
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "dist/funstream_darwin_x86_64" ./cmd/funstream/funstream.go       # Darwin 64bit
env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o "dist/funstream_windows_i386.exe" ./cmd/funstream/funstream.go     # Windows 32bit
env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "dist/funstream_windows_x86_64.exe" ./cmd/funstream/funstream.go # Windows 64bit