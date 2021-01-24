install:
	get-oui && cat ieee-oui.txt | grep -v '#' | sort > mac-fab.txt && rm ieee-oui.txt

build: mac-fab.txt
	@echo "Compiling for Linux and Windows amd64"
	GOOS=linux GOARCH=amd64 go build -o bin/go-fingerprint main_linux.go
	GOOS=windows GOARCH=amd64 go build -o bin/go-fingerprint.exe main_windows.go
	sudo setcap CAP_NET_RAW+ep bin/go-fingerprint

clean:	
	rm -rf ./bin
	go clean