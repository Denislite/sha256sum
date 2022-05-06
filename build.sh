#!/usr/bin/env bash

if [ "$1" = "win32" ]
then
echo "Building for Windows 32-bit"
env GOOS=windows GOARCH=386 go build -o hasher_32bit.exe cmd/main.go

elif [ "$1" = "win64"  ]; then
echo "Building for Windows 64-bit"
env GOOS=windows GOARCH=amd64 go build -o hasher_64bit.exe cmd/main.go

elif [ "$1" = "osx32" ]; then
echo "Building for MacOS 32-bit"
env GOOS=darwin GOARCH=386 go build -o hasher_32bit cmd/main.go

elif [ "$1" = "osx64" ]; then
echo "Building for MacOS 64-bit"
env GOOS=darwin GOARCH=amd64 go build -o hasher_64bit cmd/main.go

elif [ "$1" = "linux32" ]; then
echo "Building for Linux 32-bit"
env GOOS=linux GOARCH=386 go build -o hasher_32bit cmd/main.go

else
echo "Building for Linux 64-bit"
env GOOS=linux GOARCH=amd64 go build -o hasher_64bit cmd/main.go

fi