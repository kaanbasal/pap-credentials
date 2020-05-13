delete-all:
	@echo "Deleting binaries"
	@rm -f build/*
.PHONY: delete-all

build-darwin-amd64:
	@echo "Generating darwin:amd64"
	@rm -f build/darwin_x64
	@env GOOS=darwin GOARCH=amd64 go build -o build/darwin_x64
.PHONY: build-darwin-amd64

build-linux-amd64:
	@echo "Generating linux:amd64"
	@rm -f build/linux_x64
	@env GOOS=linux GOARCH=amd64 go build -o build/linux_x64
.PHONY: build-linux-amd64

build-windows-amd64:
	@echo "Generating windows:amd64"
	@rm -f build/windows_x64.exe
	@env GOOS=windows GOARCH=amd64 go build -o build/windows_x64.exe
.PHONY: build-windows-amd64

build-windows-386:
	@echo "Generating windows:386"
	@rm -f build/windows_x86.exe
	@env GOOS=windows GOARCH=386 go build -o build/windows_x86.exe
.PHONY: build-windows-386
