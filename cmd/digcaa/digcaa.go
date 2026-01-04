package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/weppos/digcaa"
)

func main() {
	// Custom usage banner
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "digcaa - DNS CAA Lookup Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "  digcaa <hostname>\n\n")
		fmt.Fprintf(os.Stderr, "Description:\n\n")
		fmt.Fprintf(os.Stderr, "  Fetch DNS CAA (Certification Authority Authorization) records for a hostname.\n")
		fmt.Fprintf(os.Stderr, "  The tool walks the hostname up to the root and collects all CAA records.\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n\n")
		fmt.Fprintf(os.Stderr, "  digcaa example.com\n")
		fmt.Fprintf(os.Stderr, "  digcaa www.example.com\n")
	}

	flag.Parse()

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
