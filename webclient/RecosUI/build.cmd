xcopy ..\..\general\assets .\public\assets /d /s /v /e /y /q
call npm run build
rd /s /q ..\..\service\web\webclient
xcopy .\dist\ ..\..\service\web\webclient\ /d /s /v /e /y /q
