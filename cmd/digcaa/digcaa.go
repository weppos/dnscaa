package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/weppos/dnscaa"
)

func init() {
	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	hostname := flag.Arg(0)
	fmt.Println("Lookup", hostname)
	dig(hostname)
}

func dig(hostname string) {
	records, err := dnscaa.Lookup(hostname)

	if err != nil {
		fmt.Println(err)
	}

	for _, record := range records {
		fmt.Println(record)
	}
}
