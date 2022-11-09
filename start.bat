IF [%1] == [] exit /B
IF [%2] == [] exit /B
start /D %1  go run %2
