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
	if r.resolver != DefaultResolver {
		t.Errorf("expected resolver to be %v, got %v", DefaultResolver, r.resolver)
	}
}

func TestNewResolverWithConfig(t *testing.T) {
	customTimeout := 10 * time.Second
	customResolver := "1.1.1.1:53"
	config := &Config{
		Timeout:  customTimeout,
		Resolver: customResolver,
	}
	r := NewResolverWithConfig(config)
	if r.timeout != customTimeout {
		t.Errorf("expected timeout to be %v, got %v", customTimeout, r.timeout)
	}
	if r.resolver != customResolver {
		t.Errorf("expected resolver to be %v, got %v", customResolver, r.resolver)
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
	config := &Config{
		Timeout:  customTimeout,
		Resolver: DefaultResolver,
	}
	r := NewResolverWithConfig(config)
	if r.Timeout() != customTimeout {
		t.Errorf("expected Timeout() to return %v, got %v", customTimeout, r.Timeout())
	}
}

func TestResolverGetter(t *testing.T) {
	customResolver := "1.1.1.1:53"
	config := &Config{
		Timeout:  DefaultTimeout,
		Resolver: customResolver,
	}
	r := NewResolverWithConfig(config)
	if r.Resolver() != customResolver {
		t.Errorf("expected Resolver() to return %v, got %v", customResolver, r.Resolver())
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.Timeout != DefaultTimeout {
		t.Errorf("expected config.Timeout to be %v, got %v", DefaultTimeout, config.Timeout)
	}
	if config.Resolver != DefaultResolver {
		t.Errorf("expected config.Resolver to be %v, got %v", DefaultResolver, config.Resolver)
	}
}

func TestDefaultResolver(t *testing.T) {
	if DefaultResolver != "8.8.8.8:53" {
		t.Errorf("expected DefaultResolver to be 8.8.8.8:53, got %v", DefaultResolver)
	}
}
