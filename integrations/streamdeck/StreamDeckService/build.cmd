rem "E:\Sprachen\ide\Microsoft Visual Studio\2019\Community\MSBuild\Current\Bin\MSBuild.exe" .\StreamDeckService.sln
cls
rem dotnet build .\StreamDeckService.sln --force --configuration Release --nologo
dotnet publish .\StreamDeckService.sln --force --configuration Release --nologo --no-self-contained -o bin\Release\net5.0-windows\publish\
rem -p:PublishTrimmed=true  
rem -r win10-x64