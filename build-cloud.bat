@echo off
go generate

cd ./cloud
go build -ldflags="-s -w " -o multipara.exe

move ./multipara.exe ../multipara.exe