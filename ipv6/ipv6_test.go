package ipv6

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"sync"
	"testing"
)

func TestIPv6(t *testing.T) {
	Convey("nil", t, func() {
		ips, _ := ToIPv6(nil)
		So(ips, ShouldHaveLength, 0)
	})
	Convey("error", t, func() {
		ips, err := ToIPv6([]byte("aa"))
		So(err, ShouldNotBeNil)
		So(ips, ShouldHaveLength, 0)
	})
	Convey("from to IPv6", t, func() {
		in := []byte("1234567812345678")
		ips, err := ToIPv6(in)
		So(err, ShouldBeNil)
		So(ips, ShouldHaveLength, 1)
		So(ips[0].String(), ShouldEqual, "3132:3334:3536:3738:3132:3334:3536:3738")

		out, err := FromIPv6(ips)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, in)
	})
	Convey("from to IPv6", t, func() {
		in := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		ips, err := ToIPv6(in)
		So(err, ShouldBeNil)
		So(ips, ShouldHaveLength, 2)
		So(ips[0].String(), ShouldEqual, "6161:6161:6161:6161:6161:6161:6161:6161")
		So(ips[1].String(), ShouldEqual, "6161:6161:6161:6161:6161:6161:6161:6161")

		out, err := FromIPv6(ips)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, in)
	})
	Convey("from domain", t, func() {
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
		buf, err := FromDomain("zh.test.optool.net", f)
		So(err, ShouldBeNil)
		So(buf, ShouldResemble, []byte("从前有座山の里有座庙12"))

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
		So(err1, ShouldBeNil)
		So(buf1, ShouldResemble, []byte("从前有座山の里有座庙12"))
	})
}
