@echo off
go version
echo building service
go build -ldflags="-s -w" -o recos-serice.exe cmd/service.go