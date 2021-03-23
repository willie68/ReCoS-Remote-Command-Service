xcopy ..\..\general\assets .\public\assets /s /v /e /y
call npm run build
xcopy /S /V /E /Y dist\* ..\..\service\devdata\webadmin\