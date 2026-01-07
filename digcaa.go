package digcaa

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

const (
	// DefaultTimeout is the default timeout for DNS queries.
	DefaultTimeout = 5 * time.Second
	// DefaultResolver is the default DNS resolver address.
	DefaultResolver = "8.8.8.8:53"
)

var defaultResolver = NewResolver()

// Config represents configuration options for a DNS resolver.
type Config struct {
	// Timeout is the timeout for DNS queries.
	Timeout time.Duration
	// Resolver is the DNS resolver address (e.g., "8.8.8.8:53").
	Resolver string
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Timeout:  DefaultTimeout,
		Resolver: DefaultResolver,
	}
}

// Resolver represents a DNS resolver that can be used to lookup the CAA records.
type Resolver struct {
	dnsClient *dns.Client
	timeout   time.Duration
	resolver  string
}

// NewResolver constructs a new DNS resolver with an underlying DNS client
// configured with the default timeout and resolver.
func NewResolver() *Resolver {
	return NewResolverWithConfig(DefaultConfig())
}

// NewResolverWithConfig constructs a new DNS resolver with an underlying DNS client
// configured with the specified configuration.
func NewResolverWithConfig(config *Config) *Resolver {
	r := new(Resolver)
	r.timeout = config.Timeout
	r.resolver = config.Resolver
	r.dnsClient = &dns.Client{
		Timeout: config.Timeout,
	}
	return r
}

// Timeout returns the configured timeout for DNS queries.
func (r *Resolver) Timeout() time.Duration {
	return r.timeout
}

// Resolver returns the configured DNS resolver address.
func (r *Resolver) Resolver() string {
	return r.resolver
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
	results, err := r.LookupWithQuery(hostname)
	records := make([]*dns.CAA, 0, len(results))
	for _, result := range results {
		records = append(records, result.Record)
	}
	return records, err
}

// LookupResult represents a CAA record along with the hostname that was queried.
type LookupResult struct {
	Query  string
	Record *dns.CAA
}

// LookupWithQuery performs a lookup of the CAA records for the hostname,
// and returns results annotated with the query name.
func (r *Resolver) LookupWithQuery(hostname string) ([]LookupResult, error) {
	labels := strings.Split(hostname, ".")

	var wg sync.WaitGroup
	ch := make(chan LookupResult)
	errCh := make(chan error, len(labels))

	for i := range labels {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()

			caas, err := r.LookupCAA(query)
			if err != nil {
				errCh <- err
				return
			}

			for _, caa := range caas {
				ch <- LookupResult{
					Query:  query,
					Record: caa,
				}
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

	var results []LookupResult
	for result := range ch {
		results = append(results, result)
	}

	// Wait for error collection to complete
	<-done

	if len(errs) > 0 {
		return results, errs[0]
	}

	return results, nil
}

// LookupCAA performs a DNS query to lookup the CAA records for the given hostname,
// and returns the array of records. The query will timeout according to the configured
// timeout duration (default: 5 seconds).
func (r *Resolver) LookupCAA(name string) ([]*dns.CAA, error) {
	var rrs []*dns.CAA

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeCAA)

	rsp, _, err := r.dnsClient.Exchange(msg, r.resolver)
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
