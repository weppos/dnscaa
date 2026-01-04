# digcaa - DNS CAA inspection CLI and (Go) library

The _digcaa_ package provides a command-line tool (and eventually a library) to fetch DNS CAA (Certification Authority Authorization) records for a hostname.

## CLI usage

### Compile

```
$ go build cmd/digcaa/digcaa.go
```

### Run

```
$ ./digcaa www.comodo.com

2 records found
comodo.com.	1199	IN	CAA	0 iodef "mailto:sslabuse@comodoca.com"
comodo.com.	1199	IN	CAA	0 issue "comodoca.com"
```

Alternatively, use `go run` to compile and run the command:

```
$ go run cmd/digcaa/digcaa.go www.comodo.com

2 records found
comodo.com.	1199	IN	CAA	0 iodef "mailto:sslabuse@comodoca.com"
comodo.com.	1199	IN	CAA	0 issue "comodoca.com"
```
