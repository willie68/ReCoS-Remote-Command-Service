@echo off
go version
echo building service
go build -ldflags="-s -w" -o pl-example.exe cmd/service.go
