# dnscaa - CAA DNS Library for Go

The _dnscaa_ package provides a library (not yet) and command-line tool to fetch the DNS CAA (Certification Authority Authorization) records for an hostname.

## CLI usage

#### Compile

```
$ go build cmd/digcaa/digcaa.go
```

#### Run

```
$ ./digcaa www.comodo.com

2 records found
comodo.com.	1199	IN	CAA	0 iodef "mailto:sslabuse@comodoca.com"
comodo.com.	1199	IN	CAA	0 issue "comodoca.com"
```

Or use `go run` to compile-and-run the command:

```
$ go run cmd/digcaa/digcaa.go www.comodo.com

2 records found
comodo.com.	1199	IN	CAA	0 iodef "mailto:sslabuse@comodoca.com"
comodo.com.	1199	IN	CAA	0 issue "comodoca.com"
```
