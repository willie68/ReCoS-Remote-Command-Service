@echo off

echo build service binaries
cd service
call deployments\build.cmd
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
