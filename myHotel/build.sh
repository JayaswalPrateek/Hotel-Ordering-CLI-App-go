#!/bin/bash
cd src
echo "Generating executable"
go build -o myHotel -ldflags="-s -w" main.go &&
mv myHotel ./../ &&
cd .. &&
echo "Hashing executable" &&
sha512sum myHotel > myHotel.sha512sum &&
mkdir -p out
mv myHotel out/
mv myHotel.sha512sum out/
echo "Done !"
