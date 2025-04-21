@echo off
setlocal enabledelayedexpansion

REM Define the starting port number
set /a port=8081

REM Define the number of instances to start
set instances=2

REM Loop to start each instance
for /l %%i in (1,1,%instances%) do (
    start "" cmd /k "server.exe -port !port!"
    set /a port+=1
)

echo Started %instances% instances of server.exe
pause
