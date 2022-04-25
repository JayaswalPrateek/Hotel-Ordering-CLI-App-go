#!/bin/bash
echo "Building Standard Binary"
go build -o standardBin main.go &&
echo "Building Compressed Binary"
go build -o compressedBin -ldflags="-s -w" main.go &&
echo "now creating checksum" &&
sha512sum standardBin > standardBin.sha512sum &&
mkdir -p checksum && mv standardBin.sha512sum checksum/ &&
mkdir -p exec && mv standardBin exec/ &&
sha512sum compressedBin > compressedBin.sha512sum &&
mkdir -p checksum && mv compressedBin.sha512sum checksum/ &&
mkdir -p exec && mv compressedBin exec/ &&
echo "Done !"

