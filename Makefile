run:
	go run main.go

build:
	go build

compile:
	echo "Compiling for Linux and Windows x86_64"
	GOOS=linux GOARCH=amd64 go build -o bin/go-fingerprint
	GOOS=windows GOARCH=amd64 go build -o bin/go-fingerprint.exe

clean:	
	rm -rf ./bin
	go clean