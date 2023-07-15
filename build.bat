@echo off
go generate
go build -ldflags="-s -w "