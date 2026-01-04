package digcaa

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

var defaultResolver = NewResolver()

// Resolver represents a DNS resolver that can be used to lookup the CAA records.
type Resolver struct {
	dnsClient *dns.Client
}

// NewResolver constructs a new DNS resolver with an underlying DNS client.
func NewResolver() *Resolver {
	r := new(Resolver)
	r.dnsClient = &dns.Client{}
	return r
}

// Lookup performs a DNS CAA lookup for the hostname using the default Resolver.
// See Resolver.Lookup() for more information.
func Lookup(hostname string) ([]*dns.CAA, error) {
	return defaultResolver.Lookup(hostname)
}

// Lookup performs a lookup of the CAA records for the hostname.
//
// When fetching the CAA records for an hostname, the client must walk the hostname up to the
// root and collect all the CAA records for each names that compose the hostname.
//
// This method executes all the relevant DNS queries, and returns a single array of all the CAA
// records associated to the given hostname (including parent names).
//
// For instance, the result of a Lookup() for www.example.com may contain records for
// - www.example.com
// - example.com
// - com
//
// The DNS queries are executed in parallel to minimize the execution time. Lookup() returns
// when all the DNS queries are completed.
func (r *Resolver) Lookup(hostname string) ([]*dns.CAA, error) {
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

			caas, err := r.LookupCAA(name)
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

// LookupCAA performs a DNS query to lookup the CAA records for the given hostname,
// and returns the array of records.
func (r *Resolver) LookupCAA(name string) ([]*dns.CAA, error) {
	var rrs []*dns.CAA

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeCAA)

	rsp, _, err := r.dnsClient.Exchange(msg, "8.8.8.8:53")
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
