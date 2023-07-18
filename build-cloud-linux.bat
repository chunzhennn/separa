@echo off
go generate
SET GOARCH=amd64
SET GOOS=linux
cd ./cloud

go build -ldflags="-s -w " -o multipara

move ./multipara ../multipara