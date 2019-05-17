package ipv6

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

func ToIPv6(in []byte) (ips []net.IP, err error) {
	if len(in)%net.IPv6len != 0 {
		return nil, errors.New("must be n * 16 length")
	}
	ipCount := len(in) / net.IPv6len
	ips = make([]net.IP, ipCount)
	for i := 0; i < ipCount; i++ {
		ips[i] = make(net.IP, net.IPv6len)
		copy(ips[i], in[i*net.IPv6len:(i+1)*net.IPv6len])
	}
	return
}

func FromIPv6(ips []net.IP) (out []byte, err error) {
	ipCount := len(ips)
	out = make([]byte, ipCount*net.IPv6len)
	for i := 0; i < ipCount; i++ {
		copy(out[i*net.IPv6len:(i+1)*net.IPv6len], ips[i])
	}

	return
}

func FromDomain(domain string) (out []byte, err error) {
	concurrentNum := 5
	retryCount := 3

	allIPv6 := make([]net.IP, 0, 4)

	var ipsErr error
	var ipsArray [][]net.IP
	wg := new(sync.WaitGroup)

	for i := 0; ; i++ {
		// Concurrent by group
		var successCount int32

		ipsArray = make([][]net.IP, concurrentNum, concurrentNum)

		wg.Add(concurrentNum)

		for j := 0; j < concurrentNum; j++ {
			go func(i, j int) {
				defer wg.Done()

				index := i*concurrentNum + j
				for a := 0; a < retryCount; a++ {
					ips, err := net.LookupIP(fmt.Sprintf("%02d.%s", index, domain))
					if err == nil {
						ipsArray[j] = ips
						atomic.AddInt32(&successCount, 1)

						break
					} else {
						if j == 0 {
							ipsErr = err
						}
					}
				}

			}(i, j)
		}

		wg.Wait()

		for i, ips := range ipsArray {

			if int32(i) < successCount {
				if len(ips) == 0 {
					return nil, errors.New("empty IP list")
				}
				if len(ips[0]) != net.IPv6len {
					return nil, errors.Errorf("unexpected IP: %s", ips[0])
				}
				allIPv6 = append(allIPv6, ips[0])
			}
		}

		if len(allIPv6) != 0 {

			if successCount < int32(concurrentNum) {
				break
			}
		} else {
			return nil, ipsErr
		}
	}

	out, err = FromIPv6(allIPv6)
	if err != nil {
		return nil, errors.Errorf("convert from IPv6 failed: %v", err)
	}
	return
}
