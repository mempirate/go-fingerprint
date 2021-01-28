BINPATH = ${GOPATH}

all: clean install

update-dev:
	sudo update-ieee-data
	get-oui; grep -v "#" ieee-oui.txt | sort > mac-fab.txt; rm ieee-oui.txt

update:
	@echo "[*] Getting IEEE OUI file"
	curl https://raw.githubusercontent.com/jonasbostoen/go-fingerprint/main/mac-fab.txt -o mac-fab.txt

install:
	@echo "[*] Installing..."
	go install
	@echo CAP_NET_RAW is needed to create and send packets
	sudo setcap CAP_NET_RAW+ep ${BINPATH}/bin/go-fingerprint
	@echo [+] Done

clean:	
	@echo "[*] Cleaning..."
	rm -rf ${BINPATH}/bin/go-fingerprint
	go clean
