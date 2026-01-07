# digcaa - DNS CAA inspection CLI and (Go) library

A command-line tool and Go library to fetch DNS CAA (Certification Authority Authorization) records for a hostname.

## CLI Installation

```shell
go install github.com/weppos/digcaa/cmd/digcaa@latest
```

or build from source:

```shell
git clone https://github.com/weppos/digcaa.git
cd digcaa
make build
```

The binary will be created in the current directory as `digcaa`.

## CLI Usage

```shell
digcaa <hostname>
```

Lookup CAA records for a domain:

```shell
digcaa www.dnsimple.com

# 10 records found
# www.dnsimple.com : dnsimple.com.   3600    IN      CAA     0 issue "amazon.com"
# www.dnsimple.com : dnsimple.com.   3600    IN      CAA     0 issue "letsencrypt.org"
# www.dnsimple.com : dnsimple.com.   3600    IN      CAA     0 issuewild "sectigo.com"
# www.dnsimple.com : dnsimple.com.   3600    IN      CAA     0 iodef "mailto:ops@dnsimple.com"
# [...]
```

### Configuration Options

Configure DNS query timeout and resolver:

```shell
# Use a 10 second timeout
digcaa --timeout 10s example.com

# Use a 10 millisecond timeout
digcaa --timeout 10ms example.com

# Use a 1 minute timeout
digcaa --timeout 1m example.com

# Use Cloudflare DNS resolver
digcaa --resolver 1.1.1.1:53 example.com

# Combine timeout and resolver
digcaa --timeout 10s --resolver 1.1.1.1:53 example.com
```

## Library Usage

Install the library:

```shell
go get github.com/weppos/digcaa
```

Use it in your Go code:

```go
package main

import (
	"fmt"
	"log"

	"github.com/weppos/digcaa"
)

func main() {
	// Use default resolver and timeout
	records, err := digcaa.Lookup("www.comodo.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d records found\n", len(records))
	for _, record := range records {
		fmt.Println(record)
	}
}
```

### Custom Configuration

Use custom timeout and resolver:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/weppos/digcaa"
)

func main() {
	// Create a custom configuration
	config := &digcaa.Config{
		Timeout:  10 * time.Second,
		Resolver: "1.1.1.1:53", // Cloudflare DNS
	}

	resolver := digcaa.NewResolverWithConfig(config)
	records, err := resolver.Lookup("www.example.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d records found\n", len(records))
	for _, record := range records {
		fmt.Println(record)
	}
}
```

## License

[MIT License](LICENSE.txt).
