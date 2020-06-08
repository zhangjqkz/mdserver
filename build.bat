:: prepare
rd /s /q "markdown-server"
md "markdown-server"

:: build mac
SET GOOS=darwin
SET GOARCH=amd64
go build  -o markdown-server/markdown-darwin .

:: build linux
SET GOOS=linux
SET GOARCH=amd64
go build -o markdown-server/markdown-linux .

:: build windows
SET GOOS=windows
SET GOARCH=amd64
go build -o markdown-server/markdown-windows.exe .

xcopy lib markdown-server\markdown-lib /i/y/e