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
	dig(hostname)
}

func dig(hostname string) {
	records, err := digcaa.Lookup(hostname)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d records found\n", len(records))

	for _, record := range records {
		fmt.Println(record)
	}
}
