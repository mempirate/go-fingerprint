update:
	@echo "[+] Updating IEEE OUI data"
	get-oui && cat ieee-oui.txt | grep -v '#' | sort > mac-fab.txt && rm ieee-oui.txt

download:
	@echo "[+] Getting IEEE OUI file"
	curl https://raw.githubusercontent.com/jonasbostoen/go-fingerprint/main/mac-fab.txt -o mac-fab.txt

build: mac-fab.txt
	@echo "[+] Compiling for Linux and Windows 64-bit"
	GOOS=linux GOARCH=amd64 go build -o bin/go-fingerprint
	GOOS=windows GOARCH=amd64 go build -o bin/go-fingerprint.exe
	sudo setcap CAP_NET_RAW+ep bin/go-fingerprint

install:

clean:	
	rm -rf ./bin
	go clean

all: download build