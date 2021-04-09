xcopy ..\..\general\assets .\public\assets /d /s /v /e /y
call npm run build
xcopy .\dist\ ..\..\service\web\webclient\ /d /s /v /e /y
