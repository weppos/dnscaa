.PHONY: build clean test

build:
	go build -o digcaa cmd/digcaa/digcaa.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o digcaa.exe cmd/digcaa/digcaa.go

clean:
	rm -f digcaa
