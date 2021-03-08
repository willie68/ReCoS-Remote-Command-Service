@echo off
go version
echo building service
go build -ldflags="-s -w" -o recos-service.exe cmd/service.go