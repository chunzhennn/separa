@echo off
go generate

cd ./cloud
go build -ldflags="-s -w " -o mutipara.exe

move ./mutipara.exe ../mutipara.exe