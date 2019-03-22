package ipv6

import (
	"errors"
	"net"
)

func ToIPv6(in []byte) (ips [] net.IP, err error) {
	if len(in)%net.IPv6len != 0 {
		return nil, errors.New("must be n * 16 length")
	}
	ipCount := len(in) / net.IPv6len
	ips = make([]net.IP, ipCount)
	for i := 0; i < ipCount; i ++ {
		ips[i] = make(net.IP, net.IPv6len)
		copy(ips[i], in[i*net.IPv6len:(i+1)*net.IPv6len])
	}
	return
}

func FromIPv6(ips []net.IP) (out []byte, err error) {
	ipCount := len(ips)
	out = make([]byte, ipCount * net.IPv6len)
	for i := 0; i < ipCount; i ++ {
		copy(out[i*net.IPv6len:(i+1)*net.IPv6len], ips[i])
	}

	return
}