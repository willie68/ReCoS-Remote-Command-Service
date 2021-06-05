xcopy ..\..\general\assets .\public\assets /d /s /v /e /y
call npm run build
rd ..\..\service\web\webadmin
pause
xcopy .\dist\ ..\..\service\web\webadmin\ /d /s /v /e /y
