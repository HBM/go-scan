
.PHONY: statik
statik:
	statik -src=public

dev:
	go run main.go 
linux:
	go build -o scan

win64:
	GOOS=windows GOARCH=amd64 go build -o scan64.exe

win32:
	GOOS=windows GOARCH=386 go build -o scan32.exe

mac:
	GOOS=darwin GOARCH=amd64 go build -o scan_mac

release: linux win64 win32 mac
