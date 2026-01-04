# digcaa - DNS CAA inspection CLI and (Go) library

A command-line tool and Go library to fetch DNS CAA (Certification Authority Authorization) records for a hostname.

## CLI Installation

```shell
go install github.com/weppos/dnscaa/cmd/digcaa@latest
```

or build from source:

```shell
git clone https://github.com/weppos/dnscaa.git
cd dnscaa
make build
```

The binary will be created in the current directory as `digcaa`.

## CLI Usage

```shell
digcaa <hostname>
```

Lookup CAA records for a domain:

```shell
digcaa www.comodo.com

# 2 records found
# comodo.com.	1199	IN	CAA	0 iodef "mailto:sslabuse@comodoca.com"
# comodo.com.	1199	IN	CAA	0 issue "comodoca.com"
```

## Library Usage

Install the library:

```shell
go get github.com/weppos/dnscaa
```

Use it in your Go code:

```go
package main

import (
	"fmt"
	"log"

	"github.com/weppos/dnscaa"
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
