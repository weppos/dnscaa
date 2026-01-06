package digcaa

import (
	"fmt"
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
//
// If any DNS query fails, the first error encountered is returned along with any successfully
// retrieved records.
func (r *Resolver) Lookup(hostname string) ([]*dns.CAA, error) {
	var records []*dns.CAA
	labels := strings.Split(hostname, ".")

	var wg sync.WaitGroup
	ch := make(chan *dns.CAA)
	errCh := make(chan error, len(labels))

	for i := range labels {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			caas, err := r.LookupCAA(name)
			if err != nil {
				errCh <- err
				return
			}

			for _, caa := range caas {
				ch <- caa
			}
		}(strings.Join(labels[i:], "."))
	}

	go func() {
		wg.Wait()
		close(ch)
		close(errCh)
	}()

	// Collect errors in a separate goroutine
	var errs []error
	done := make(chan struct{})
	go func() {
		for err := range errCh {
			errs = append(errs, err)
		}
		close(done)
	}()

	for rr := range ch {
		records = append(records, rr)
	}

	// Wait for error collection to complete
	<-done

	if len(errs) > 0 {
		return records, errs[0]
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
		return nil, fmt.Errorf("CAA lookup failed for %s: %w", name, err)
	}

	if rsp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("CAA lookup for %s returned %s", name, dns.RcodeToString[rsp.Rcode])
	}

	for _, rr := range rsp.Answer {
		if cca, ok := rr.(*dns.CAA); ok {
			rrs = append(rrs, cca)
		}
	}

	return rrs, nil
}
