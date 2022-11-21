@echo off
go version
echo building service
go build -ldflags="-s -w" -o recos-service-console.exe cmd/service/main.go
