package ipv6

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"testing"
)

func TestIPv6(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ips, _ := ToIPv6(nil)
		if len(ips) != 0 {
			t.Errorf("ToIPv6(nil) = %v, want no IPs", ips)
		}
	})
	t.Run("error", func(t *testing.T) {
		ips, err := ToIPv6([]byte("aa"))
		if err == nil || len(ips) != 0 {
			t.Errorf("ToIPv6(aa) = %v, %v; want an error", ips, err)
		}
	})
	t.Run("from to IPv6", func(t *testing.T) {
		in := []byte("1234567812345678")
		ips, err := ToIPv6(in)
		if err != nil || len(ips) != 1 || ips[0].String() != "3132:3334:3536:3738:3132:3334:3536:3738" {
			t.Fatalf("ToIPv6(%s) = %v, %v; want [3132:3334:3536:3738:3132:3334:3536:3738]", in, ips, err)
		}

		out, err := FromIPv6(ips)
		if err != nil || !bytes.Equal(out, in) {
			t.Errorf("FromIPv6(%v) = %q, %v; want %q", ips, out, err, in)
		}
	})
	t.Run("from to IPv6", func(t *testing.T) {
		in := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		ips, err := ToIPv6(in)
		if err != nil || len(ips) != 2 {
			t.Fatalf("ToIPv6(%s) = %v, %v; want 2 IPs", in, ips, err)
		}
		for _, ip := range ips {
			if ip.String() != "6161:6161:6161:6161:6161:6161:6161:6161" {
				t.Errorf("ToIPv6(%s) returned %s, want 6161:6161:6161:6161:6161:6161:6161:6161", in, ip)
			}
		}

		out, err := FromIPv6(ips)
		if err != nil || !bytes.Equal(out, in) {
			t.Errorf("FromIPv6(%v) = %q, %v; want %q", ips, out, err, in)
		}
	})
	t.Run("from domain", func(t *testing.T) {
		// The AAAA records of the README example.
		records := map[string]net.IP{
			"00.zh.test.optool.net": net.ParseIP("e4bb:8ee5:898d:e69c:89e5:baa7:e5b1:b1e3"),
			"01.zh.test.optool.net": net.ParseIP("81ae:e987:8ce6:9c89:e5ba:a7e5:ba99:3132"),
		}
		f := func(host string) ([]net.IP, error) {
			if ip, ok := records[host]; ok {
				return []net.IP{ip}, nil
			}
			return nil, errors.New("no such host")
		}
		want := []byte("从前有座山の里有座庙12")
		buf, err := FromDomain("zh.test.optool.net", f)
		if err != nil || !bytes.Equal(buf, want) {
			t.Errorf("FromDomain() = %q, %v; want %q", buf, err, want)
		}

		// Retry: each name fails its first 2 lookups, and FromDomain tries 3 times.
		var mu sync.Mutex
		lookups := make(map[string]int)
		f1 := func(host string) ([]net.IP, error) {
			mu.Lock()
			lookups[host]++
			n := lookups[host]
			mu.Unlock()
			if n <= 2 {
				return nil, errors.New("no such host")
			}
			return f(host)
		}
		buf1, err1 := FromDomain("zh.test.optool.net", f1)
		if err1 != nil || !bytes.Equal(buf1, want) {
			t.Errorf("FromDomain() with failing lookups = %q, %v; want %q", buf1, err1, want)
		}
	})
}
