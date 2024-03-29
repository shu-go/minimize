set NAME=minimize

minimize

set GOOS=windows
set GOARCH=386
set ZIPNAME=%NAME%_%GOOS%_%GOARCH%.zip 
go build -ldflags "-s -w -H=windowsgui" %*
del %ZIPNAME% 2>nul
zip %ZIPNAME% *.exe README.txt LICENSE.txt

set GOOS=windows
set GOARCH=amd64
set ZIPNAME=%NAME%_%GOOS%_%GOARCH%.zip 
rem go build -ldflags "-s -w -H=windowsgui" %*
rem del %ZIPNAME% 2>nul
rem zip %ZIPNAME% *.exe README.txt LICENSE.txt

minimize -r
