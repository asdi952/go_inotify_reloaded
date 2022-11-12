@echo off
IF [%1] == [] exit /B
IF [%2] == [] exit /B

cmd /k "cd %1 && go run %2"