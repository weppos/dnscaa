package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

var client *dns.Client

func main() {
	hostname := os.Args[1]
	client = new(dns.Client)

	fmt.Println("Lookup", hostname)
	records, err := Lookup(hostname)

	if err != nil {
		fmt.Println(err)
	}

	for _, record := range records {
		fmt.Println(record)
	}
}

func Lookup(hostname string) ([]*dns.CAA, error) {
	var records []*dns.CAA
	labels := strings.Split(hostname, ".")

	var wg sync.WaitGroup
	ch := make(chan *dns.CAA, 1)

	// Current issues:
	// - no error checking
	// - it should probably quit on first error
	// See https://blog.golang.org/pipelines (select + range)
	for i := range labels {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			caas, err := LookupCAA(name)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, caa := range caas {
				caa := caa
				ch <- caa
			}
		}(strings.Join(labels[i:], "."))
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for rr := range ch {
		records = append(records, rr)
	}

	return records, nil
}

func LookupCAA(name string) ([]*dns.CAA, error) {
	var rrs []*dns.CAA

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeCAA)

	rsp, err := dns.Exchange(msg, "8.8.8.8:53")
	if err != nil {
		log.Println("CAA lookup failed", name, err)
		return nil, err
	}

	if rsp.Rcode != dns.RcodeSuccess {
		log.Println("CAA lookup not success", name, dns.RcodeToString[rsp.Rcode])
		return nil, fmt.Errorf("lookup code %s", dns.RcodeToString[rsp.Rcode])
	}

	for _, rr := range rsp.Answer {
		if cca, ok := rr.(*dns.CAA); ok {
			rrs = append(rrs, cca)
		}
	}

	return rrs, nil
}
