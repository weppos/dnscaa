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
	resolver := flag.String("resolver", digcaa.DefaultResolver, "DNS resolver address (e.g., 8.8.8.8:53, 1.1.1.1:53)")

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
		fmt.Fprintf(os.Stderr, "  digcaa --resolver 1.1.1.1:53 example.com\n")
		fmt.Fprintf(os.Stderr, "  digcaa --timeout 10s --resolver 1.1.1.1:53 example.com\n")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	hostname := flag.Arg(0)
	dig(hostname, *timeout, *resolver)
}

func dig(hostname string, timeout time.Duration, resolverAddr string) {
	config := &digcaa.Config{
		Timeout:  timeout,
		Resolver: resolverAddr,
	}
	resolver := digcaa.NewResolverWithConfig(config)
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
