#!/bin/bash
echo "Generating executable"
go build -o myHotel main.go &&
echo "Hashing executable" &&
sha512sum myHotel > myHotel.sha512sum &&
echo "Done !"
