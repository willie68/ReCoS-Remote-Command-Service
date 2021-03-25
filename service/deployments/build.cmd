@echo off
go version
echo building service
go build -ldflags="-s -w -H=windowsgui" -o recos-service.exe cmd/service.go