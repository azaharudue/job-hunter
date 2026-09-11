@echo off
REM Daily job-hunt run. Invoked by Windows Task Scheduler. Usage: job-hunt-daily.bat
set PROJECT=%~dp0..
set LOG=%PROJECT%\joblists\job-hunt-run.log
if not exist "%PROJECT%\joblists" mkdir "%PROJECT%\joblists"
echo [%date% %time%] start >> "%LOG%"
cd /d "%PROJECT%"
opencode run --agent job-hunter --auto "Run your daily job-hunt routine now." >> "%LOG%" 2>&1
echo [%date% %time%] exit=%errorlevel% >> "%LOG%"