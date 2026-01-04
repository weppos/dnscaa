.PHONY: build clean test

build:
	go build -o digcaa cmd/digcaa/digcaa.go

clean:
	rm -f digcaa
