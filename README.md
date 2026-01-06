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
# dnsimple.com.   3600    IN      CAA     0 issue "amazon.com"
# dnsimple.com.   3600    IN      CAA     0 issue "letsencrypt.org"
# dnsimple.com.   3600    IN      CAA     0 issuewild "sectigo.com"
# dnsimple.com.   3600    IN      CAA     0 iodef "mailto:ops@dnsimple.com"
# [...]
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

## License

[MIT License](LICENSE.txt).
