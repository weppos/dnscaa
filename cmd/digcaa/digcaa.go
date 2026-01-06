package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/weppos/digcaa"
)

func main() {
	// Define flags
	timeout := flag.Duration("timeout", digcaa.DefaultTimeout, "Timeout for DNS queries (e.g., 5s, 10s, 1m)")

	// Custom usage banner
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "digcaa - DNS CAA Lookup Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "  digcaa [options] <hostname>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDescription:\n\n")
		fmt.Fprintf(os.Stderr, "  Fetch DNS CAA (Certification Authority Authorization) records for a hostname.\n")
		fmt.Fprintf(os.Stderr, "  The tool walks the hostname up to the root and collects all CAA records.\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n\n")
		fmt.Fprintf(os.Stderr, "  digcaa example.com\n")
		fmt.Fprintf(os.Stderr, "  digcaa www.example.com\n")
		fmt.Fprintf(os.Stderr, "  digcaa --timeout 10s example.com\n")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	hostname := flag.Arg(0)
	dig(hostname, *timeout)
}

func dig(hostname string, timeout time.Duration) {
	resolver := digcaa.NewResolverWithTimeout(timeout)
	records, err := resolver.Lookup(hostname)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d records found\n", len(records))

	for _, record := range records {
		fmt.Println(record)
	}
}
