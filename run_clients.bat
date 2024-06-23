@echo off
set NUM_CLIENTS=30  REM Number of clients to run

for /L %%i in (1,1,%NUM_CLIENTS%) do (
    start /B "" cmd /C client.exe
)

echo Started %NUM_CLIENTS% clients.
pause
