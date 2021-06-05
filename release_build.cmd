@echo off
".\3rd party\GoVersionSetter.exe" -i -m
".\3rd party\GoVersionSetter.exe" -e npm -f ./webclient/recosadmin/package.json
".\3rd party\GoVersionSetter.exe" -e npm -f ./webclient/RecosUI/package.json
".\3rd party\GoVersionSetter.exe" -e iss -f ./install/setup.iss -o MyAppVersion
".\3rd party\GoVersionSetter.exe" -e vs -f ./integrations/streamdeck/StreamDeckService/StreamDeckService.csproj
".\3rd party\GoVersionSetter.exe" -e gores -f ./service/winres/winres.json -o RT_MANIFEST/#1/0409/identity/version,RT_VERSION/#1/0000/fixed/file_version,RT_VERSION/#1/0000/fixed/product_version,RT_VERSION/#1/0000/info/0409/ProductVersion,RT_VERSION/#1/0000/info/0409/FileVersion

call sync_assets.cmd

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

echo build streamdeck integration
cd integrations\streamdeck\StreamDeckService
call build.cmd
cd ..
cd ..
cd ..

echo build service binaries
cd service
call deployments\build.cmd
deployments\go-winres.exe make
deployments\go-winres.exe patch recos-service.exe
cd ..

echo build setup
cd install
iscc setup.iss
cd ..

echo release ready. please test it.
pause