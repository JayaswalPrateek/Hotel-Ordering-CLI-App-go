#!/bin/bash

go version &&
cd src &&
mkdir -p ./../out &&
echo -e "\n\n"


echo "Checking updates for dependencies" &&
go get -u && go mod tidy &&
echo -e "DONE\n\n"


echo "Linting" &&
go fmt ./main.go && go vet ./main.go && golangci-lint run ./main.go &&
echo -e "DONE\n\n"


echo "Generating executable" &&
go build -a -ldflags="-s -w" -o ./../out/myHotel main.go &&
cd ./../out/ &&
echo -e "DONE\n\n"


echo "Generating checksum (SHA512)" &&
sha512sum myHotel > .myHotel.sha512sum &&
echo "DONE"
