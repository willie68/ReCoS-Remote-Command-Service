@echo off

echo build service binaries
cd service
call deployments\build.cmd
deployments\go-winres.exe make
deployments\go-winres.exe patch recos-service.exe
cd ..

echo build web client
cd webclient\RecosUI
call build.cmd
cd ..
cd ..

echo build web admin
cd webclient\recosadmin
call build.cmd
cd ..
cd ..

echo build setup
cd install
iscc setup.iss
cd ..
