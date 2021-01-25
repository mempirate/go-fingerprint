BINPATH = ${GOPATH}

update:
	@echo "[+] Getting IEEE OUI file"
	curl https://raw.githubusercontent.com/jonasbostoen/go-fingerprint/main/mac-fab.txt -o mac-fab.txt

install:
	@echo "[+] Installing..."
	go install
	# CAP_NET_RAW is needed to create and send packets
	sudo setcap CAP_NET_RAW+ep ${BINPATH}/bin/go-fingerprint

clean:	
	@echo "[+] Cleaning..."
	rm -rf ${BINPATH}/bin/go-fingerprint
	go clean

all: clean install