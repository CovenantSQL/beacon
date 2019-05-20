package ipv6

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"net"
	"testing"
	"time"
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
		f := func(host string) ([]net.IP, error) {
			return net.LookupIP(host)
		}
		buf, err := FromDomain("zh.test.optool.net", f)
		So(err, ShouldBeNil)
		So(buf, ShouldResemble, []byte("从前有座山の里有座庙12"))

		// Retry when parsing IP errors
		rand.Seed(time.Now().UnixNano())
		f1 := func(host string) ([]net.IP, error) {
			t := rand.Intn(10000)
			if t < 2001 {
				return nil, errors.New("no such host")
			}
			return net.LookupIP(host)
		}
		buf1, err1 := FromDomain("zh.test.optool.net", f1)
		So(err1, ShouldBeNil)
		So(buf1, ShouldResemble, []byte("从前有座山の里有座庙12"))
	})
}
