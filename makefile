.PHONY: all

all: windows-386 windows-amd64 linux-386 linux-amd64 wasm

windows-386:
	GOOS=windows GOARCH=386 go build -v -o boofutils-windows-386.exe .

windows-amd64:
	GOOS=windows GOARCH=amd64 go build -v -o boofutils-windows-amd64.exe .

linux-386:
	GOOS=linux GOARCH=386 go build -v -o boofutils-linux-386 .

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -v -o boofutils-linux-amd64 .

wasm:
	GOOS=js GOARCH=wasm go build -v -o boofutils-wasm.wasm .