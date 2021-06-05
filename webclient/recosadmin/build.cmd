xcopy ..\..\general\assets .\public\assets /d /s /v /e /y /q
call npm run build
rd /s /q ..\..\service\web\webadmin
xcopy .\dist\ ..\..\service\web\webadmin\ /d /s /v /e /y /q
