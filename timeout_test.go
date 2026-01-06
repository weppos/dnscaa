package digcaa

import (
	"testing"
	"time"
)

func TestNewResolver(t *testing.T) {
	r := NewResolver()
	if r.timeout != DefaultTimeout {
		t.Errorf("expected timeout to be %v, got %v", DefaultTimeout, r.timeout)
	}
	if r.dnsClient.Timeout != DefaultTimeout {
		t.Errorf("expected client timeout to be %v, got %v", DefaultTimeout, r.dnsClient.Timeout)
	}
}

func TestNewResolverWithTimeout(t *testing.T) {
	customTimeout := 10 * time.Second
	r := NewResolverWithTimeout(customTimeout)
	if r.timeout != customTimeout {
		t.Errorf("expected timeout to be %v, got %v", customTimeout, r.timeout)
	}
	if r.dnsClient.Timeout != customTimeout {
		t.Errorf("expected client timeout to be %v, got %v", customTimeout, r.dnsClient.Timeout)
	}
}

func TestDefaultTimeout(t *testing.T) {
	if DefaultTimeout != 5*time.Second {
		t.Errorf("expected DefaultTimeout to be 5s, got %v", DefaultTimeout)
	}
}

func TestResolverTimeout(t *testing.T) {
	customTimeout := 3 * time.Second
	r := NewResolverWithTimeout(customTimeout)
	if r.Timeout() != customTimeout {
		t.Errorf("expected Timeout() to return %v, got %v", customTimeout, r.Timeout())
	}
}
